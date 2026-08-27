package templates

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

//go:embed assets/*
var EmbeddedFS embed.FS

// Extract extracts embedded assets to the target destination directory.
func Extract(subDir, targetDir string, overwrite bool) error {
	entries, err := fs.ReadDir(EmbeddedFS, "assets/"+subDir)
	if err != nil {
		return fmt.Errorf("failed to read embedded dir assets/%s: %w", subDir, err)
	}

	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return fmt.Errorf("failed to create target dir %s: %w", targetDir, err)
	}

	for _, entry := range entries {
		srcPath := "assets/" + subDir + "/" + entry.Name()
		dstPath := filepath.Join(targetDir, entry.Name())

		if entry.IsDir() {
			if err := Extract(subDir+"/"+entry.Name(), dstPath, overwrite); err != nil {
				return err
			}
			continue
		}

		if !overwrite && fileExists(dstPath) {
			continue
		}

		content, err := EmbeddedFS.ReadFile(srcPath)
		if err != nil {
			return fmt.Errorf("failed to read embedded file %s: %w", srcPath, err)
		}

		if err := os.WriteFile(dstPath, content, 0644); err != nil {
			return fmt.Errorf("failed to write %s: %w", dstPath, err)
		}
	}

	return nil
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
