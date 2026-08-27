package backup

import (
	"archive/tar"
	"compress/gzip"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/neflalabs/c4ignite/internal/config"
)

type Metadata struct {
	CreatedAt string `json:"created_at"`
	Hostname  string `json:"hostname"`
	User      string `json:"user"`
	Version   string `json:"version"`
	Encrypted bool   `json:"encrypted"`
	FileCount int    `json:"file_count"`
}

// Create creates a tar.gz backup of the src/ directory with metadata.
func Create(pCtx *config.ProjectContext, name string, encryptKey string) (string, error) {
	if err := os.MkdirAll(pCtx.BackupDir, 0755); err != nil {
		return "", err
	}

	timestamp := time.Now().Format("20060102-150405")
	if name == "" {
		name = fmt.Sprintf("c4ignite-backup-%s.tar.gz", timestamp)
	} else if !strings.HasSuffix(name, ".tar.gz") {
		name = name + ".tar.gz"
	}

	destPath := filepath.Join(pCtx.BackupDir, name)
	tempArchive := destPath + ".tmp"

	outFile, err := os.Create(tempArchive)
	if err != nil {
		return "", err
	}
	defer func() {
		outFile.Close()
		os.Remove(tempArchive)
	}()

	gw := gzip.NewWriter(outFile)
	tw := tar.NewWriter(gw)

	hostname, _ := os.Hostname()
	user := os.Getenv("USER")
	if user == "" {
		user = "unknown"
	}

	fileCount := 0
	err = filepath.Walk(pCtx.SrcPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(pCtx.SrcPath, path)
		if err != nil {
			return err
		}

		if relPath == "." {
			return nil
		}

		// Skip vendor and writable cache/logs to keep backups lean
		if strings.HasPrefix(relPath, "vendor") ||
			strings.HasPrefix(relPath, "writable/cache") ||
			strings.HasPrefix(relPath, "writable/logs") ||
			strings.HasPrefix(relPath, "writable/session") {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		header, err := tar.FileInfoHeader(info, info.Name())
		if err != nil {
			return err
		}
		header.Name = filepath.ToSlash(relPath)

		if err := tw.WriteHeader(header); err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		f, err := os.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()

		_, err = io.Copy(tw, f)
		fileCount++
		return err
	})

	if err != nil {
		return "", fmt.Errorf("failed to archive src: %w", err)
	}

	// Write metadata.json into the archive
	meta := Metadata{
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
		Hostname:  hostname,
		User:      user,
		Version:   "1.0.0",
		Encrypted: encryptKey != "",
		FileCount: fileCount,
	}

	metaBytes, _ := json.MarshalIndent(meta, "", "  ")
	metaHeader := &tar.Header{
		Name:    ".c4ignite-metadata.json",
		Size:    int64(len(metaBytes)),
		Mode:    0644,
		ModTime: time.Now(),
	}
	if err := tw.WriteHeader(metaHeader); err == nil {
		tw.Write(metaBytes)
	}

	tw.Close()
	gw.Close()
	outFile.Close()

	if encryptKey != "" {
		if err := encryptFile(tempArchive, destPath, encryptKey); err != nil {
			return "", fmt.Errorf("encryption failed: %w", err)
		}
	} else {
		if err := os.Rename(tempArchive, destPath); err != nil {
			return "", err
		}
	}

	return destPath, nil
}

// Restore extracts a tar.gz backup into the src/ directory.
func Restore(pCtx *config.ProjectContext, archivePath string, decryptKey string) error {
	var srcReader io.Reader
	f, err := os.Open(archivePath)
	if err != nil {
		return err
	}
	defer f.Close()
	srcReader = f

	if decryptKey != "" {
		tempDecrypted := archivePath + ".decrypted.tmp"
		if err := decryptFile(archivePath, tempDecrypted, decryptKey); err != nil {
			return fmt.Errorf("decryption failed: %w", err)
		}
		defer os.Remove(tempDecrypted)
		df, err := os.Open(tempDecrypted)
		if err != nil {
			return err
		}
		defer df.Close()
		srcReader = df
	}

	gr, err := gzip.NewReader(srcReader)
	if err != nil {
		return fmt.Errorf("invalid gzip archive: %w", err)
	}
	defer gr.Close()

	tr := tar.NewReader(gr)
	if err := os.MkdirAll(pCtx.SrcPath, 0755); err != nil {
		return err
	}

	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		if hdr.Name == ".c4ignite-metadata.json" {
			continue
		}

		target := filepath.Join(pCtx.SrcPath, hdr.Name)
		switch hdr.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(target, 0755); err != nil {
				return err
			}
		case tar.TypeReg:
			if err := os.MkdirAll(filepath.Dir(target), 0755); err != nil {
				return err
			}
			outF, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR|os.O_TRUNC, os.FileMode(hdr.Mode))
			if err != nil {
				return err
			}
			_, err = io.Copy(outF, tr)
			outF.Close()
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func encryptFile(src, dst, key string) error {
	data, err := os.ReadFile(src)
	if err != nil {
		return err
	}

	// Pad or hash key to 32 bytes for AES-256
	keyBytes := make([]byte, 32)
	copy(keyBytes, []byte(key))

	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return err
	}

	ciphertext := gcm.Seal(nonce, nonce, data, nil)
	return os.WriteFile(dst, ciphertext, 0644)
}

func decryptFile(src, dst, key string) error {
	ciphertext, err := os.ReadFile(src)
	if err != nil {
		return err
	}

	keyBytes := make([]byte, 32)
	copy(keyBytes, []byte(key))

	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return err
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return fmt.Errorf("ciphertext too short")
	}

	nonce, encryptedData := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, encryptedData, nil)
	if err != nil {
		return fmt.Errorf("invalid password or corrupted archive: %w", err)
	}

	return os.WriteFile(dst, plaintext, 0644)
}
