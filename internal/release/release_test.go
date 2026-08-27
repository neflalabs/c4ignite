package release

import (
	"context"
	"testing"

	"github.com/neflalabs/c4ignite/internal/config"
)

func TestBuildReleasePipelineSteps(t *testing.T) {
	pCtx := &config.ProjectContext{RootPath: "/project"}

	// Default full pipeline
	opts := ReleaseOptions{}
	steps := BuildReleasePipeline(pCtx, nil, opts)
	if len(steps) != 4 {
		t.Fatalf("expected 4 release steps, got %d", len(steps))
	}

	// Skip migration and health
	optsSkip := ReleaseOptions{
		SkipMigration: true,
		SkipHealth:    true,
	}
	stepsSkip := BuildReleasePipeline(pCtx, nil, optsSkip)
	if len(stepsSkip) != 2 {
		t.Fatalf("expected 2 release steps when skipping, got %d", len(stepsSkip))
	}
}

func TestExecutePipelineSuccess(t *testing.T) {
	pCtx := &config.ProjectContext{RootPath: "/project"}
	opts := ReleaseOptions{
		SkipMigration: true,
		SkipHealth:    true,
	}

	err := Execute(context.Background(), pCtx, nil, opts)
	if err != nil {
		t.Fatalf("expected pipeline to succeed, got %v", err)
	}
}
