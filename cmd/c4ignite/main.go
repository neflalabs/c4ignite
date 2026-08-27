package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/neflalabs/c4ignite/internal/backup"
	"github.com/neflalabs/c4ignite/internal/compose"
	"github.com/neflalabs/c4ignite/internal/config"
	"github.com/neflalabs/c4ignite/internal/doctor"
	"github.com/neflalabs/c4ignite/internal/env"
	"github.com/neflalabs/c4ignite/internal/templates"
)

var (
	Version = "dev"
)

const (
	DefaultRelease = "v4.6.3"
	AppStarterURL  = "https://github.com/codeigniter4/appstarter/archive/refs/tags/%s.tar.gz"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		return
	}

	command := os.Args[1]
	args := os.Args[2:]

	// Context for graceful signal cancellation
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	switch command {
	case "version", "--version", "-v":
		fmt.Printf("c4ignite %s (go)\n", Version)
	case "help", "--help", "-h":
		printUsage()
	case "init":
		handleInit(ctx, args)
	case "doctor":
		handleDoctor(ctx)
	case "up", "u":
		handleUp(ctx, args)
	case "down", "d", "stop":
		handleDown(ctx, args)
	case "restart":
		handleRestart(ctx, args)
	case "status", "ps":
		handleStatus(ctx)
	case "logs":
		handleLogs(ctx, args)
	case "pull":
		handlePull(ctx)
	case "shell":
		handleShell(ctx, args)
	case "spark":
		handleSpark(ctx, args)
	case "composer":
		handleComposer(ctx, args)
	case "php":
		handlePHP(ctx, args)
	case "db", "mysql":
		handleDB(ctx, args)
	case "migrate":
		handleMigrate(ctx, args)
	case "seed":
		handleSeed(ctx, args)
	case "test":
		handleTest(ctx, args)
	case "lint":
		handleLint(ctx, args)
	case "xdebug":
		handleXdebug(ctx, args)
	case "backup":
		handleBackup(ctx, args)
	case "update", "self-update":
		handleSelfUpdate(ctx)
	case "completion":
		handleCompletion(args)
	default:
		fmt.Fprintf(os.Stderr, "Error: Unknown command '%s'\n\n", command)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	banner := `
   ______ __ __  _               _  __       
  / ____// // / (_)____ _ ____  (_)/ /_ ___  
 / /    / // /_/ // __ '// __ \/ // __// _ \ 
/ /___ /__  __// // /_/ // / / / // /_ /  __/ 
\____/   /_/  /_/ \__, //_/ /_/_/ \__/ \___/  
                 /____/  CodeIgniter 4 Toolkit (Go)
`
	fmt.Println(banner)
	fmt.Printf("Usage: c4ignite <command> [arguments]\n\n")
	fmt.Println("🚀 Core Commands:")
	fmt.Println("  init [ProjectName] [--force]        Bootstrap fresh CI4 app (Default: Codeigniter4)")
	fmt.Println("  up [--build] [-d] [--with=service]  Start containerized stack")
	fmt.Println("  down [-v]                           Stop and remove containers")
	fmt.Println("  restart [service]                   Restart all or specific service")
	fmt.Println("  status                              Show live container health")
	fmt.Println("  logs [-f] [service]                 Tail service container logs")
	fmt.Println("  pull                                Pull latest container images")
	fmt.Println()
	fmt.Println("⚡ CodeIgniter 4 Ergonomics:")
	fmt.Println("  spark [command]                     Run CodeIgniter 4 Spark command")
	fmt.Println("  migrate                             Shortcut for 'spark migrate'")
	fmt.Println("  seed [seeder]                       Shortcut for 'spark db:seed'")
	fmt.Println("  composer [command]                  Execute Composer in PHP container")
	fmt.Println("  php [script]                        Execute PHP script in container")
	fmt.Println("  db [options]                        Open interactive MySQL/MariaDB shell")
	fmt.Println("  shell [service]                     Drop into bash/sh shell inside container")
	fmt.Println("  xdebug [on|off|status]              Toggle Xdebug dynamically")
	fmt.Println()
	fmt.Println("🛠️ Utilities & Testing:")
	fmt.Println("  doctor                              Run host diagnostic checks")
	fmt.Println("  update                              Self-update c4ignite to latest release")
	fmt.Println("  test [options]                      Run PHPUnit test suite")
	fmt.Println("  lint                                Run PHP code style linter")
	fmt.Println("  backup [create|restore]             Fast native backup & restore")
	fmt.Println("  completion [install|uninstall]      Configure shell auto-completions")
	fmt.Println("  version                             Show CLI version")
	fmt.Println()
}

func handleSelfUpdate(ctx context.Context) {
	fmt.Println("🔄 Checking and updating c4ignite to the latest release...")
	cmd := exec.CommandContext(ctx, "bash", "-c", "curl -fsSL https://raw.githubusercontent.com/neflalabs/c4ignite/main/install.sh | bash")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Update failed: %v\n", err)
		os.Exit(1)
	}
}

func getRunner() (*config.ProjectContext, *compose.Runner) {
	pCtx, err := config.FindProjectRoot("")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	runner, err := compose.NewRunner(pCtx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	return pCtx, runner
}

func handleInit(ctx context.Context, args []string) {
	tag := ""
	force := false
	noInstall := false
	projectName := ""

	for _, a := range args {
		if strings.HasPrefix(a, "--version=") {
			tag = strings.TrimPrefix(a, "--version=")
		} else if strings.HasPrefix(a, "--dir=") {
			projectName = strings.TrimPrefix(a, "--dir=")
		} else if strings.HasPrefix(a, "-d=") {
			projectName = strings.TrimPrefix(a, "-d=")
		} else if a == "--force" || a == "-f" {
			force = true
		} else if a == "--no-install" {
			noInstall = true
		} else if !strings.HasPrefix(a, "-") && projectName == "" {
			projectName = a
		}
	}

	// Interactive prompt if projectName not specified via flags/args
	if projectName == "" {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("📂 Project name [default: Codeigniter4]: ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		if input != "" {
			projectName = input
		} else {
			projectName = "Codeigniter4"
		}
	}

	// Resolve latest version dynamically from GitHub if not specified
	if tag == "" || tag == "latest" {
		fmt.Print("🔍 Checking latest CodeIgniter 4 release from GitHub... ")
		resolvedTag := resolveLatestCI4Release(ctx)
		tag = resolvedTag
		fmt.Printf("(%s)\n", tag)
	}

	cwd, _ := os.Getwd()
	targetDir := filepath.Join(cwd, projectName)

	if !force && dirHasFiles(targetDir) {
		fmt.Printf("⚠️  Directory '%s/' already exists and is not empty. Use --force to overwrite.\n", projectName)
		return
	}

	fmt.Printf("📦 Initializing CodeIgniter 4 (%s) into '%s/'...\n", tag, projectName)
	tarURL := fmt.Sprintf(AppStarterURL, tag)
	cacheDir := getAppCacheDir()
	_ = os.MkdirAll(cacheDir, 0755)
	tarPath := filepath.Join(cacheDir, fmt.Sprintf("appstarter-%s.tar.gz", tag))

	if !fileExists(tarPath) || force {
		fmt.Printf("⬇️  Downloading AppStarter from %s...\n", tarURL)
		if err := downloadFile(tarURL, tarPath); err != nil {
			fmt.Fprintf(os.Stderr, "Download failed: %v\n", err)
			os.Exit(1)
		}
	} else {
		fmt.Printf("⚡ Using cached archive %s\n", tarPath)
	}

	fmt.Printf("📂 Scaffolding project into %s/...\n", projectName)
	if err := extractTarGzStripped(tarPath, targetDir); err != nil {
		fmt.Fprintf(os.Stderr, "Extraction failed: %v\n", err)
		os.Exit(1)
	}

	// Extract embedded .c4ignite configurations directly into project
	_ = templates.Extract(".c4ignite", filepath.Join(targetDir, ".c4ignite"), false)
	_ = templates.Extract("env", filepath.Join(targetDir, ".c4ignite", "env"), false)

	// Copy dev.env to targetDir/.env if not present
	srcEnv := filepath.Join(targetDir, ".env")
	if !fileExists(srcEnv) {
		devEnv := filepath.Join(targetDir, ".c4ignite", "env", "dev.env")
		if fileExists(devEnv) {
			copyFile(devEnv, srcEnv)
		}
	}

	// Run composer install to provision vendor/ and framework Boot.php
	if !noInstall {
		fmt.Println("📦 Installing framework dependencies (composer install)...")
		// Temporarily change directory to targetDir to run composer install
		origDir, _ := os.Getwd()
		_ = os.Chdir(targetDir)
		projCtx, runner := getRunner()
		if runner != nil && projCtx != nil {
			if err := runner.ExecPHP(ctx, "composer", "install", "--no-interaction", "--prefer-dist"); err != nil {
				fmt.Println("⚠️  Composer install skipped. You can run 'c4ignite composer install' later.")
			} else {
				fmt.Println("✅ Framework dependencies installed successfully.")
			}
		}
		_ = os.Chdir(origDir)
	}

	fmt.Printf("\n✨ CodeIgniter 4 project '%s' created successfully!\n\n", projectName)
	fmt.Println("👉 Next steps:")
	fmt.Printf("  cd %s\n", projectName)
	fmt.Println("  c4ignite up")
	fmt.Println()
}

// resolveLatestCI4Release queries GitHub API for the latest release tag
func resolveLatestCI4Release(ctx context.Context) string {
	apiURL := "https://api.github.com/repos/codeigniter4/appstarter/releases/latest"
	reqCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(reqCtx, "GET", apiURL, nil)
	if err != nil {
		return DefaultRelease
	}
	req.Header.Set("User-Agent", "c4ignite-cli")

	client := &http.Client{Timeout: 3 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return DefaultRelease
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return DefaultRelease
	}

	var payload struct {
		TagName string `json:"tag_name"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil || payload.TagName == "" {
		return DefaultRelease
	}

	return payload.TagName
}

func handleDoctor(ctx context.Context) {
	pCtx, _ := config.FindProjectRoot("")
	fmt.Println("🩺 Running c4ignite environment diagnostics...")
	fmt.Println()
	results := doctor.RunDiagnostics(pCtx)
	allPassed := true
	for _, r := range results {
		status := "✅ PASS"
		if !r.Passed {
			status = "❌ FAIL"
			allPassed = false
		}
		fmt.Printf("  %-25s %s — %s\n", r.Name, status, r.Message)
	}
	fmt.Println()
	if allPassed {
		fmt.Println("🎉 All checks passed! Your environment is ready.")
	} else {
		fmt.Println("⚠️  Some checks failed. Please review the output above.")
	}
}

func handleUp(ctx context.Context, args []string) {
	_, runner := getRunner()
	detach := true
	build := false
	var profiles []string

	for _, a := range args {
		if a == "--build" {
			build = true
		}
		if a == "--attach" || a == "-a" {
			detach = false
		}
		if strings.HasPrefix(a, "--with=") {
			profiles = append(profiles, strings.Split(strings.TrimPrefix(a, "--with="), ",")...)
		}
	}

	fmt.Println("🚀 Spinning up c4ignite stack...")
	if err := runner.Up(ctx, detach, build, profiles); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	if detach {
		fmt.Println("✨ Services started! Web accessible at: http://localhost:8000")
	}
}

func handleDown(ctx context.Context, args []string) {
	_, runner := getRunner()
	removeVols := false
	for _, a := range args {
		if a == "-v" || a == "--volumes" {
			removeVols = true
		}
	}
	fmt.Println("🛑 Stopping c4ignite stack...")
	if err := runner.Down(ctx, removeVols); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func handleRestart(ctx context.Context, args []string) {
	_, runner := getRunner()
	svc := ""
	if len(args) > 0 {
		svc = args[0]
	}
	if err := runner.Restart(ctx, svc); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func handleStatus(ctx context.Context) {
	_, runner := getRunner()
	_ = runner.Status(ctx)
}

func handleLogs(ctx context.Context, args []string) {
	_, runner := getRunner()
	follow := false
	svc := ""
	for _, a := range args {
		if a == "-f" || a == "--follow" {
			follow = true
		} else if !strings.HasPrefix(a, "-") {
			svc = a
		}
	}
	_ = runner.Logs(ctx, svc, follow)
}

func handlePull(ctx context.Context) {
	_, runner := getRunner()
	_ = runner.Pull(ctx)
}

func handleShell(ctx context.Context, args []string) {
	_, runner := getRunner()
	svc := "php"
	if len(args) > 0 {
		svc = args[0]
	}
	if err := runner.Shell(ctx, svc); err != nil {
		fmt.Fprintf(os.Stderr, "Shell error: %v\n", err)
		os.Exit(1)
	}
}

func handleSpark(ctx context.Context, args []string) {
	_, runner := getRunner()
	cmdArgs := append([]string{"php", "spark"}, args...)
	if err := runner.ExecPHP(ctx, cmdArgs...); err != nil {
		os.Exit(1)
	}
}

func handleMigrate(ctx context.Context, args []string) {
	handleSpark(ctx, append([]string{"migrate"}, args...))
}

func handleSeed(ctx context.Context, args []string) {
	handleSpark(ctx, append([]string{"db:seed"}, args...))
}

func handleComposer(ctx context.Context, args []string) {
	_, runner := getRunner()
	cmdArgs := append([]string{"composer"}, args...)
	if err := runner.ExecPHP(ctx, cmdArgs...); err != nil {
		os.Exit(1)
	}
}

func handlePHP(ctx context.Context, args []string) {
	_, runner := getRunner()
	cmdArgs := append([]string{"php"}, args...)
	if err := runner.ExecPHP(ctx, cmdArgs...); err != nil {
		os.Exit(1)
	}
}

func handleDB(ctx context.Context, args []string) {
	pCtx, runner := getRunner()
	envFile := filepath.Join(pCtx.SrcPath, ".env")
	envMap, _ := env.Load(envFile)

	dbName := envMap["database.default.database"]
	if dbName == "" {
		dbName = "app"
	}
	user := envMap["database.default.username"]
	if user == "" {
		user = "app"
	}
	pass := envMap["database.default.password"]
	if pass == "" {
		pass = "secret"
	}

	if err := runner.ExecMySQL(ctx, dbName, user, pass); err != nil {
		fmt.Fprintf(os.Stderr, "Database connection error: %v\n", err)
		os.Exit(1)
	}
}

func handleTest(ctx context.Context, args []string) {
	_, runner := getRunner()
	cmdArgs := append([]string{"vendor/bin/phpunit"}, args...)
	if err := runner.ExecPHP(ctx, cmdArgs...); err != nil {
		os.Exit(1)
	}
}

func handleLint(ctx context.Context, args []string) {
	_, runner := getRunner()
	cmdArgs := []string{"vendor/bin/phpcs", "app"}
	if err := runner.ExecPHP(ctx, cmdArgs...); err != nil {
		fmt.Fprintf(os.Stderr, "Lint error: %v\n", err)
		os.Exit(1)
	}
}

func handleXdebug(ctx context.Context, args []string) {
	if len(args) == 0 {
		fmt.Println("Usage: c4ignite xdebug [on|off|status]")
		return
	}
	pCtx, runner := getRunner()
	iniPath := filepath.Join(pCtx.DockerDevDir, "php", "conf.d", "xdebug.ini")
	switch args[0] {
	case "on":
		_ = os.WriteFile(iniPath, []byte("xdebug.mode=debug,develop\nxdebug.start_with_request=yes\n"), 0644)
		_ = runner.Restart(ctx, "php")
		fmt.Println("🐞 Xdebug enabled (mode=debug,develop)")
	case "off":
		_ = os.WriteFile(iniPath, []byte("xdebug.mode=off\n"), 0644)
		_ = runner.Restart(ctx, "php")
		fmt.Println("⚡ Xdebug disabled (mode=off)")
	case "status":
		content, _ := os.ReadFile(iniPath)
		fmt.Printf("Xdebug configuration (%s):\n%s\n", iniPath, string(content))
	}
}

func handleBackup(ctx context.Context, args []string) {
	if len(args) == 0 {
		fmt.Println("Usage: c4ignite backup [create|restore] [options]")
		return
	}
	pCtx, _ := getRunner()
	switch args[0] {
	case "create":
		name := ""
		key := ""
		if len(args) > 1 && !strings.HasPrefix(args[1], "-") {
			name = args[1]
		}
		for _, a := range args {
			if strings.HasPrefix(a, "--key=") {
				key = strings.TrimPrefix(a, "--key=")
			}
		}
		path, err := backup.Create(pCtx, name, key)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Backup failed: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("📦 Backup created: %s\n", path)
	case "restore":
		if len(args) < 2 {
			fmt.Println("Usage: c4ignite backup restore <archive-path> [--key=secret]")
			return
		}
		archivePath := args[1]
		key := ""
		for _, a := range args {
			if strings.HasPrefix(a, "--key=") {
				key = strings.TrimPrefix(a, "--key=")
			}
		}
		if err := backup.Restore(pCtx, archivePath, key); err != nil {
			fmt.Fprintf(os.Stderr, "Restore failed: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("✅ Backup restored successfully into src/")
	}
}

func handleCompletion(args []string) {
	if len(args) > 0 && (args[0] == "install" || args[0] == "setup") {
		installCompletion()
		return
	}
	if len(args) > 0 && (args[0] == "uninstall" || args[0] == "remove") {
		uninstallCompletion()
		return
	}

	shell := "bash"
	if len(args) > 0 {
		shell = args[0]
	}
	switch shell {
	case "bash":
		fmt.Println(`_c4ignite() {
    local cur prev words cword
    _init_completion -n : 2>/dev/null || {
        cur="${COMP_WORDS[COMP_CWORD]}"
        prev="${COMP_WORDS[COMP_CWORD-1]}"
        words=("${COMP_WORDS[@]}")
        cword=$COMP_CWORD
    }

    local commands="up down restart status logs pull shell spark composer php db migrate seed test lint xdebug backup doctor update init version help completion"

    if [ "$cword" -eq 1 ]; then
        COMPREPLY=( $(compgen -W "${commands}" -- "$cur") )
        return 0
    fi

    local command="${words[1]}"
    case "$command" in
        completion)
            COMPREPLY=( $(compgen -W "bash zsh fish install uninstall" -- "$cur") )
            ;;
        xdebug)
            COMPREPLY=( $(compgen -W "on off status" -- "$cur") )
            ;;
        backup)
            COMPREPLY=( $(compgen -W "create restore" -- "$cur") )
            ;;
        init)
            COMPREPLY=( $(compgen -W "--force --version= --dir=" -- "$cur") )
            ;;
        up)
            COMPREPLY=( $(compgen -W "--build --with= -d" -- "$cur") )
            ;;
        down)
            COMPREPLY=( $(compgen -W "-v --volumes" -- "$cur") )
            ;;
        shell|restart|logs)
            COMPREPLY=( $(compgen -W "php nginx mysql" -- "$cur") )
            ;;
    esac
    return 0
}
complete -F _c4ignite c4ignite`)
	case "zsh":
		fmt.Println(`if ! type compdef >/dev/null 2>&1; then
    autoload -Uz compinit && compinit -u
fi

_c4ignite() {
    local -a commands
    commands=(
        'up:Start containerized stack'
        'down:Stop and remove containers'
        'restart:Restart services'
        'status:Show container health'
        'logs:Tail service logs'
        'pull:Pull latest container images'
        'shell:Drop into bash/sh shell inside container'
        'spark:Execute CI4 Spark command'
        'db:Open interactive MySQL console'
        'migrate:Run database migrations (spark migrate)'
        'seed:Seed database (spark db:seed)'
        'composer:Run Composer in PHP container'
        'php:Run PHP script in container'
        'test:Run PHPUnit test suite'
        'lint:Run PHP code style linter'
        'xdebug:Toggle Xdebug dynamically (on|off|status)'
        'doctor:Run environment diagnostic checks'
        'update:Self-update c4ignite to latest release'
        'backup:Backup or restore application files'
        'completion:Generate or install shell autocompletions'
        'init:Bootstrap fresh CodeIgniter 4 project'
        'version:Show version information'
        'help:Show help information'
    )

    if (( CURRENT == 2 )); then
        _describe -t commands 'c4ignite command' commands
        return
    fi

    local command="$words[2]"
    case "$command" in
        completion)
            local -a comp_cmds
            comp_cmds=('bash:Bash completion' 'zsh:Zsh completion' 'fish:Fish completion' 'install:Install completion to profile' 'uninstall:Remove completion from profile')
            _describe -t subcommands 'completion action' comp_cmds
            ;;
        xdebug)
            local -a xdebug_cmds
            xdebug_cmds=('on:Enable Xdebug' 'off:Disable Xdebug' 'status:Show Xdebug configuration')
            _describe -t subcommands 'xdebug mode' xdebug_cmds
            ;;
        backup)
            local -a backup_cmds
            backup_cmds=('create:Create backup archive' 'restore:Restore backup archive')
            _describe -t subcommands 'backup action' backup_cmds
            ;;
        shell|restart|logs)
            local -a services
            services=('php:PHP Application Service' 'nginx:Web Server' 'mysql:Database')
            _describe -t services 'service' services
            ;;
    esac
}

compdef _c4ignite c4ignite 2>/dev/null || true`)
	case "fish":
		fmt.Println(`complete -c c4ignite -f
complete -c c4ignite -n "__fish_use_subcommand" -a "up down restart status logs pull shell spark composer php db migrate seed test lint xdebug backup doctor init version help completion"`)
	}
}

func installCompletion() {
	shell := os.Getenv("SHELL")
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error finding home directory: %v\n", err)
		return
	}

	if strings.Contains(shell, "zsh") {
		zshrc := filepath.Join(homeDir, ".zshrc")
		appendSnippetIfMissing(zshrc, `
# c4ignite autocompletion
eval "$(c4ignite completion zsh)"
`)
		fmt.Printf("✅ Added c4ignite completion to %s. Run 'source ~/.zshrc' to activate.\n", zshrc)
	} else if strings.Contains(shell, "fish") {
		fishDir := filepath.Join(homeDir, ".config", "fish", "completions")
		_ = os.MkdirAll(fishDir, 0755)
		fishFile := filepath.Join(fishDir, "c4ignite.fish")
		_ = os.WriteFile(fishFile, []byte(`complete -c c4ignite -f
complete -c c4ignite -n "__fish_use_subcommand" -a "up down restart status logs pull shell spark composer php db migrate seed test lint xdebug backup doctor init version help completion"
`), 0644)
		fmt.Printf("✅ Installed fish completion at %s\n", fishFile)
	} else {
		// Default to bash
		bashrc := filepath.Join(homeDir, ".bashrc")
		appendSnippetIfMissing(bashrc, `
# c4ignite autocompletion
eval "$(c4ignite completion bash)"
`)
		fmt.Printf("✅ Added c4ignite completion to %s. Run 'source ~/.bashrc' to activate.\n", bashrc)
	}
}

func uninstallCompletion() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return
	}

	// Remove from ~/.bashrc
	bashrc := filepath.Join(homeDir, ".bashrc")
	removeSnippetIfExists(bashrc, "c4ignite completion")

	// Remove from ~/.zshrc
	zshrc := filepath.Join(homeDir, ".zshrc")
	removeSnippetIfExists(zshrc, "c4ignite completion")

	// Remove from ~/.config/fish/completions/c4ignite.fish
	fishFile := filepath.Join(homeDir, ".config", "fish", "completions", "c4ignite.fish")
	_ = os.Remove(fishFile)
	fmt.Println("🧹 Removed c4ignite shell completion configurations.")
}

func removeSnippetIfExists(filePath, marker string) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return
	}
	lines := strings.Split(string(data), "\n")
	var filtered []string
	skipNext := false
	for _, l := range lines {
		if strings.Contains(l, marker) || strings.Contains(l, "c4ignite autocompletion") {
			skipNext = true
			continue
		}
		if skipNext && strings.TrimSpace(l) == "" {
			skipNext = false
			continue
		}
		skipNext = false
		filtered = append(filtered, l)
	}
	_ = os.WriteFile(filePath, []byte(strings.Join(filtered, "\n")), 0644)
}

func appendSnippetIfMissing(filePath, snippet string) {
	content, _ := os.ReadFile(filePath)
	if strings.Contains(string(content), "c4ignite completion") {
		return
	}
	f, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return
	}
	defer f.Close()
	_, _ = f.WriteString(snippet)
}

func getAppCacheDir() string {
	userCache, err := os.UserCacheDir()
	if err == nil && userCache != "" {
		return filepath.Join(userCache, "c4ignite")
	}
	homeDir, err := os.UserHomeDir()
	if err == nil && homeDir != "" {
		return filepath.Join(homeDir, ".config", "c4ignite", "cache")
	}
	return filepath.Join(os.TempDir(), "c4ignite-cache")
}

// Helpers
func downloadFile(url, dst string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, resp.Body)
	return err
}

func extractTarGzStripped(src, dst string) error {
	cmd := exec.Command("tar", "-xzf", src, "--strip-components=1", "-C", dst)
	_ = os.MkdirAll(dst, 0755)
	return cmd.Run()
}

func dirHasFiles(path string) bool {
	entries, err := os.ReadDir(path)
	if err != nil {
		return false
	}
	return len(entries) > 0
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, in)
	return err
}
