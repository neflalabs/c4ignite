package builder

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/neflalabs/c4ignite/internal/config"
)

// BuildOptions holds options for container image building
type BuildOptions struct {
	Tag        string
	Target     string
	NoCache    bool
	Push       bool
	Dockerfile string
	ContextDir string
}

// BuildCommand generates the docker build command arguments
func BuildCommand(pCtx *config.ProjectContext, opts BuildOptions) []string {
	dockerfile := opts.Dockerfile
	if dockerfile == "" {
		dockerfile = filepath.Join(pCtx.RootPath, "docker", "prod", "Dockerfile")
	}

	contextDir := opts.ContextDir
	if contextDir == "" {
		contextDir = pCtx.RootPath
	}

	tag := opts.Tag
	if tag == "" {
		tag = "c4ignite-app:latest"
	}

	args := []string{"build", "-t", tag, "-f", dockerfile}

	if opts.Target != "" {
		args = append(args, "--target", opts.Target)
	}

	if opts.NoCache {
		args = append(args, "--no-cache")
	}

	args = append(args, contextDir)
	return args
}

// Build executes docker build using system docker CLI
func Build(ctx context.Context, pCtx *config.ProjectContext, opts BuildOptions) error {
	dockerfile := opts.Dockerfile
	if dockerfile == "" {
		dockerfile = filepath.Join(pCtx.RootPath, "docker", "prod", "Dockerfile")
	}

	if _, err := os.Stat(dockerfile); os.IsNotExist(err) {
		return fmt.Errorf("production Dockerfile not found at %s", dockerfile)
	}

	args := BuildCommand(pCtx, opts)
	fmt.Printf("🔨 Building production image '%s'...\n", opts.Tag)
	fmt.Printf("   Dockerfile: %s\n", dockerfile)

	cmd := exec.CommandContext(ctx, "docker", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("docker build failed: %w", err)
	}

	fmt.Printf("✅ Image built successfully: %s\n", opts.Tag)

	if opts.Push {
		fmt.Printf("🚀 Pushing image '%s' to registry...\n", opts.Tag)
		pushCmd := exec.CommandContext(ctx, "docker", "push", opts.Tag)
		pushCmd.Stdout = os.Stdout
		pushCmd.Stderr = os.Stderr
		if err := pushCmd.Run(); err != nil {
			return fmt.Errorf("docker push failed: %w", err)
		}
		fmt.Printf("✨ Image pushed successfully: %s\n", opts.Tag)
	}

	return nil
}
