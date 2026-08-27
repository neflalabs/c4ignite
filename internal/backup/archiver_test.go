package backup

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/neflalabs/c4ignite/internal/config"
)

func TestBackupCreateAndRestore(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "c4ignite-backup-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	srcDir := filepath.Join(tempDir, "src")
	_ = os.MkdirAll(srcDir, 0755)
	_ = os.WriteFile(filepath.Join(srcDir, "app.php"), []byte("<?php echo 'hello';"), 0644)

	pCtx := &config.ProjectContext{
		RootPath:  tempDir,
		SrcPath:   srcDir,
		BackupDir: filepath.Join(tempDir, "backups"),
	}

	archive, err := Create(pCtx, "test.tar.gz", "secret-key-123456")
	if err != nil {
		t.Fatalf("failed to create backup: %v", err)
	}

	// Remove source file
	_ = os.Remove(filepath.Join(srcDir, "app.php"))

	// Restore
	if err := Restore(pCtx, archive, "secret-key-123456"); err != nil {
		t.Fatalf("failed to restore backup: %v", err)
	}

	restoredContent, err := os.ReadFile(filepath.Join(srcDir, "app.php"))
	if err != nil {
		t.Fatalf("restored file missing: %v", err)
	}
	if string(restoredContent) != "<?php echo 'hello';" {
		t.Errorf("unexpected content: %s", string(restoredContent))
	}
}
