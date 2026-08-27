package env

import (
	"os"
	"path/filepath"
	"testing"
)

func TestEnvLoadAndSet(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "c4ignite-env-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	envFile := filepath.Join(tempDir, ".env")
	initialContent := "# Test Config\nAPP_ENV=development\nPORT=8000\n"
	if err := os.WriteFile(envFile, []byte(initialContent), 0644); err != nil {
		t.Fatal(err)
	}

	envMap, err := Load(envFile)
	if err != nil {
		t.Fatalf("failed to load env: %v", err)
	}

	if envMap["APP_ENV"] != "development" || envMap["PORT"] != "8000" {
		t.Errorf("unexpected env map values: %v", envMap)
	}

	if err := SetKey(envFile, "PORT", "9000"); err != nil {
		t.Fatalf("failed to set key: %v", err)
	}

	updatedMap, _ := Load(envFile)
	if updatedMap["PORT"] != "9000" {
		t.Errorf("expected PORT=9000, got %s", updatedMap["PORT"])
	}
}
