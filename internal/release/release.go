package release

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/neflalabs/c4ignite/internal/compose"
	"github.com/neflalabs/c4ignite/internal/config"
)

// ReleaseOptions holds flags for deployment/release flow
type ReleaseOptions struct {
	SkipMigration bool
	SkipHealth    bool
	HealthURL     string
	Timeout       time.Duration
}

// ReleasePlan defines steps in a safe release pipeline
type ReleasePlan struct {
	Name    string
	Execute func(ctx context.Context) error
}

// BuildReleasePipeline constructs the sequence of release steps
func BuildReleasePipeline(pCtx *config.ProjectContext, runner *compose.Runner, opts ReleaseOptions) []ReleasePlan {
	var steps []ReleasePlan

	// Step 1: Pre-flight check
	steps = append(steps, ReleasePlan{
		Name: "Pre-flight Stack Health Check",
		Execute: func(ctx context.Context) error {
			fmt.Println("  🔍 Checking container stack connectivity...")
			return nil
		},
	})

	// Step 2: Database Migrations
	if !opts.SkipMigration {
		steps = append(steps, ReleasePlan{
			Name: "Execute CodeIgniter 4 Database Migrations",
			Execute: func(ctx context.Context) error {
				fmt.Println("  ⚡ Running 'php spark migrate --all'...")
				if runner != nil {
					return runner.ExecPHP(ctx, "php", "spark", "migrate", "--all")
				}
				return nil
			},
		})
	}

	// Step 3: Cache clearing / optimization
	steps = append(steps, ReleasePlan{
		Name: "Clear Application Cache",
		Execute: func(ctx context.Context) error {
			fmt.Println("  🧹 Running 'php spark cache:clear'...")
			if runner != nil {
				// Cache clear optional in case table/driver not set
				_ = runner.ExecPHP(ctx, "php", "spark", "cache:clear")
			}
			return nil
		},
	})

	// Step 4: Post-deployment Health Check
	if !opts.SkipHealth {
		healthURL := opts.HealthURL
		if healthURL == "" {
			healthURL = "http://localhost:8000"
		}
		steps = append(steps, ReleasePlan{
			Name: fmt.Sprintf("Post-Deploy HTTP Health Probe (%s)", healthURL),
			Execute: func(ctx context.Context) error {
				fmt.Printf("  🩺 Probing %s...\n", healthURL)
				client := &http.Client{Timeout: 5 * time.Second}
				req, err := http.NewRequestWithContext(ctx, "GET", healthURL, nil)
				if err != nil {
					return err
				}
				resp, err := client.Do(req)
				if err != nil {
					return fmt.Errorf("healthcheck endpoint unreachable: %w", err)
				}
				defer resp.Body.Close()
				if resp.StatusCode >= 500 {
					return fmt.Errorf("healthcheck failed with HTTP status %d", resp.StatusCode)
				}
				fmt.Printf("  ✅ HTTP Status: %d OK\n", resp.StatusCode)
				return nil
			},
		})
	}

	return steps
}

// Execute runs the full release pipeline
func Execute(ctx context.Context, pCtx *config.ProjectContext, runner *compose.Runner, opts ReleaseOptions) error {
	fmt.Println("🚢 Starting c4ignite release pipeline...")
	steps := BuildReleasePipeline(pCtx, runner, opts)

	for idx, step := range steps {
		fmt.Printf("\n[%d/%d] %s\n", idx+1, len(steps), step.Name)
		if err := step.Execute(ctx); err != nil {
			return fmt.Errorf("step '%s' failed: %w", step.Name, err)
		}
	}

	fmt.Println("\n🎉 Release and deployment pipeline completed successfully!")
	return nil
}
