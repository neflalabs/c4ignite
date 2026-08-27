package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFindProjectRoot(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "c4ignite-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	composeDir := filepath.Join(tempDir, "docker", "dev")
	if err := os.MkdirAll(composeDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(composeDir, "docker-compose.yml"), []byte("services:\n"), 0644); err != nil {
		t.Fatal(err)
	}

	nestedDir := filepath.Join(tempDir, "src", "app", "Controllers")
	if err := os.MkdirAll(nestedDir, 0755); err != nil {
		t.Fatal(err)
	}

	ctx, err := FindProjectRoot(nestedDir)
	if err != nil {
		t.Fatalf("expected to find root, got error: %v", err)
	}

	if ctx.RootPath != tempDir {
		t.Errorf("expected root path %s, got %s", tempDir, ctx.RootPath)
	}
}

func TestResolveAppDirWithMarker(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "c4ignite-appdir-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	// Write marker file
	markerPath := filepath.Join(tempDir, ".c4ignite-app")
	if err := os.WriteFile(markerPath, []byte("my_custom_app\n"), 0644); err != nil {
		t.Fatal(err)
	}

	resolved := ResolveAppDir(tempDir)
	if resolved != "my_custom_app" {
		t.Errorf("expected 'my_custom_app', got '%s'", resolved)
	}
}
