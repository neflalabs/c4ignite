package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type ProjectContext struct {
	RootPath     string
	AppDirName   string
	SrcPath      string
	DockerDevDir string
	ComposeFile  string
	BackupDir    string
	HostUID      int
	HostGID      int
}

// FindProjectRoot traverses upward from the start directory looking for c4ignite markers.
func FindProjectRoot(startDir string) (*ProjectContext, error) {
	if startDir == "" {
		var err error
		startDir, err = os.Getwd()
		if err != nil {
			return nil, fmt.Errorf("failed to get current working directory: %w", err)
		}
	}

	curr, err := filepath.Abs(startDir)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve absolute path: %w", err)
	}

	for {
		// Check for docker/dev/docker-compose.yml or scripts/c4ignite or c4ignite.yml
		composeCheck := filepath.Join(curr, "docker", "dev", "docker-compose.yml")
		scriptCheck := filepath.Join(curr, "scripts", "c4ignite")
		configCheck := filepath.Join(curr, "c4ignite.yml")

		if fileExists(composeCheck) || fileExists(scriptCheck) || fileExists(configCheck) {
			uid := os.Getuid()
			gid := os.Getgid()
			if uid == 0 {
				uid = 1000
			}
			if gid == 0 {
				gid = 1000
			}

			// Resolve App Directory name dynamically
			appDirName := ResolveAppDir(curr)

			return &ProjectContext{
				RootPath:     curr,
				AppDirName:   appDirName,
				SrcPath:      filepath.Join(curr, appDirName),
				DockerDevDir: filepath.Join(curr, "docker", "dev"),
				ComposeFile:  composeCheck,
				BackupDir:    filepath.Join(curr, "backups"),
				HostUID:      uid,
				HostGID:      gid,
			}, nil
		}

		parent := filepath.Dir(curr)
		if parent == curr {
			break
		}
		curr = parent
	}

	return nil, fmt.Errorf("not in a c4ignite project (no docker/dev/docker-compose.yml or scripts/c4ignite found in parent directories)")
}

// ResolveAppDir determines which folder contains the CI4 application
func ResolveAppDir(rootDir string) string {
	// 1. Check environment variable
	if envDir := os.Getenv("C4IGNITE_APP_DIR"); envDir != "" {
		return envDir
	}

	// 2. Check marker file .c4ignite-app
	markerPath := filepath.Join(rootDir, ".c4ignite-app")
	if data, err := os.ReadFile(markerPath); err == nil {
		trimmed := strings.TrimSpace(string(data))
		if trimmed != "" {
			return trimmed
		}
	}

	// 3. Check if default src/ exists and has spark or composer.json
	if fileExists(filepath.Join(rootDir, "src", "spark")) || fileExists(filepath.Join(rootDir, "src", "composer.json")) {
		return "src"
	}

	// 4. Check if app/ exists and has spark or composer.json
	if fileExists(filepath.Join(rootDir, "app", "spark")) || fileExists(filepath.Join(rootDir, "app", "composer.json")) {
		return "app"
	}

	// Default fallback
	return "src"
}

func fileExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !info.IsDir()
}
