package builder

import (
	"path/filepath"
	"reflect"
	"testing"

	"github.com/neflalabs/c4ignite/internal/config"
)

func TestBuildCommand(t *testing.T) {
	pCtx := &config.ProjectContext{
		RootPath: "/project/root",
		SrcPath:  "/project/root/src",
	}

	opts := BuildOptions{
		Tag:     "myapp:v1.0.0",
		Target:  "runtime",
		NoCache: true,
	}

	expected := []string{
		"build",
		"-t", "myapp:v1.0.0",
		"-f", filepath.Join("/project/root", "docker", "prod", "Dockerfile"),
		"--target", "runtime",
		"--no-cache",
		"/project/root",
	}

	cmd := BuildCommand(pCtx, opts)
	if !reflect.DeepEqual(cmd, expected) {
		t.Fatalf("expected command %v, got %v", expected, cmd)
	}
}

func TestBuildCommandDefaults(t *testing.T) {
	pCtx := &config.ProjectContext{
		RootPath: "/project/root",
	}

	opts := BuildOptions{}
	expected := []string{
		"build",
		"-t", "c4ignite-app:latest",
		"-f", filepath.Join("/project/root", "docker", "prod", "Dockerfile"),
		"/project/root",
	}

	cmd := BuildCommand(pCtx, opts)
	if !reflect.DeepEqual(cmd, expected) {
		t.Fatalf("expected command %v, got %v", expected, cmd)
	}
}
