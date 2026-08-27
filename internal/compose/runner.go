package compose

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/neflalabs/c4ignite/internal/config"
)

type Runner struct {
	ctx        *config.ProjectContext
	composeCmd []string
}

func NewRunner(pCtx *config.ProjectContext) (*Runner, error) {
	// Detect docker compose or docker-compose
	if checkCmd("docker", "compose", "version") {
		return &Runner{
			ctx:        pCtx,
			composeCmd: []string{"docker", "compose"},
		}, nil
	} else if checkCmd("docker-compose", "version") {
		return &Runner{
			ctx:        pCtx,
			composeCmd: []string{"docker-compose"},
		}, nil
	}
	return nil, fmt.Errorf("docker compose is required but not installed or not working")
}

func checkCmd(name string, args ...string) bool {
	cmd := exec.Command(name, args...)
	return cmd.Run() == nil
}

func (r *Runner) baseCmd(ctx context.Context, additionalArgs ...string) *exec.Cmd {
	args := append([]string{}, r.composeCmd[1:]...)
	args = append(args, "-f", r.ctx.ComposeFile)
	args = append(args, additionalArgs...)

	cmd := exec.CommandContext(ctx, r.composeCmd[0], args...)
	cmd.Dir = r.ctx.RootPath
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Inject environment variables
	env := os.Environ()
	env = append(env, fmt.Sprintf("HOST_UID=%d", r.ctx.HostUID))
	env = append(env, fmt.Sprintf("HOST_GID=%d", r.ctx.HostGID))
	if r.ctx.AppDirName != "" {
		env = append(env, fmt.Sprintf("C4IGNITE_APP_DIR=%s", r.ctx.AppDirName))
	}
	cmd.Env = env
	return cmd
}

func (r *Runner) Up(ctx context.Context, detach bool, build bool, profiles []string) error {
	var args []string
	for _, p := range profiles {
		if strings.TrimSpace(p) != "" {
			args = append(args, "--profile", strings.TrimSpace(p))
		}
	}
	args = append(args, "up")
	if detach {
		args = append(args, "-d")
	}
	if build {
		args = append(args, "--build")
	}
	cmd := r.baseCmd(ctx, args...)
	return cmd.Run()
}

func (r *Runner) Down(ctx context.Context, removeVolumes bool) error {
	args := []string{"down"}
	if removeVolumes {
		args = append(args, "-v")
	}
	cmd := r.baseCmd(ctx, args...)
	return cmd.Run()
}

func (r *Runner) Restart(ctx context.Context, service string) error {
	args := []string{"restart"}
	if service != "" {
		args = append(args, service)
	}
	cmd := r.baseCmd(ctx, args...)
	return cmd.Run()
}

func (r *Runner) Status(ctx context.Context) error {
	cmd := r.baseCmd(ctx, "ps")
	return cmd.Run()
}

func (r *Runner) Logs(ctx context.Context, service string, follow bool) error {
	args := []string{"logs"}
	if follow {
		args = append(args, "-f")
	}
	if service != "" {
		args = append(args, service)
	}
	cmd := r.baseCmd(ctx, args...)
	return cmd.Run()
}

func (r *Runner) Pull(ctx context.Context) error {
	cmd := r.baseCmd(ctx, "pull")
	return cmd.Run()
}

func (r *Runner) Build(ctx context.Context, noCache bool) error {
	args := []string{"build"}
	if noCache {
		args = append(args, "--no-cache")
	}
	cmd := r.baseCmd(ctx, args...)
	return cmd.Run()
}

// ExecPHP runs a command inside the php service container with proper workspace
func (r *Runner) ExecPHP(ctx context.Context, args ...string) error {
	// First try 'run --rm' for fast reliable execution
	runArgs := []string{"run", "--rm", "--no-deps", "-w", "/var/www/html", "php"}
	runArgs = append(runArgs, args...)
	cmd := r.baseCmd(ctx, runArgs...)
	return cmd.Run()
}

// ExecMySQL opens an interactive shell to the database container with credentials pre-loaded
func (r *Runner) ExecMySQL(ctx context.Context, dbName, user, password string) error {
	if dbName == "" {
		dbName = "app"
	}
	if user == "" {
		user = "app"
	}
	if password == "" {
		password = "secret"
	}
	execArgs := []string{"exec", "mysql", "mariadb", "-u" + user, "-p" + password, dbName}
	cmd := r.baseCmd(ctx, execArgs...)
	return cmd.Run()
}

// Shell opens an interactive bash/sh shell inside the requested service
func (r *Runner) Shell(ctx context.Context, service string) error {
	if service == "" {
		service = "php"
	}
	shellBin := "bash"
	if service == "python" || service == "alpine" {
		shellBin = "sh"
	}
	execArgs := []string{"exec", service, shellBin}
	cmd := r.baseCmd(ctx, execArgs...)
	err := cmd.Run()
	if err != nil {
		// Fallback to sh if bash is not present
		fallbackArgs := []string{"exec", service, "sh"}
		fallbackCmd := r.baseCmd(ctx, fallbackArgs...)
		return fallbackCmd.Run()
	}
	return nil
}
