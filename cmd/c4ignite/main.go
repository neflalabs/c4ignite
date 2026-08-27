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
	"github.com/neflalabs/c4ignite/internal/builder"
	"github.com/neflalabs/c4ignite/internal/compose"
	"github.com/neflalabs/c4ignite/internal/config"
	"github.com/neflalabs/c4ignite/internal/doctor"
	"github.com/neflalabs/c4ignite/internal/env"
	"github.com/neflalabs/c4ignite/internal/release"
	"github.com/neflalabs/c4ignite/internal/templates"
)

const (
	Version        = "v2026.08.02"
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
		fmt.Printf("c4ignite %s (native golang)\n", Version)
	case "help", "--help", "-h":
		printUsage()
	case "init":
		handleInit(ctx, args)
	case "doctor":
		handleDoctor(ctx)
	case "up":
		handleUp(ctx, args)
	case "down":
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
	case "build":
		handleBuild(ctx, args)
	case "deploy", "release":
		handleRelease(ctx, args)
	case "backup":
		handleBackup(ctx, args)
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
                 /____/  CodeIgniter 4 Toolkit (Native Go)
`
	fmt.Println(banner)
	fmt.Printf("Usage: c4ignite <command> [arguments]\n\n")
	fmt.Println("🚀 Core Lifecycle Commands:")
	fmt.Println("  up [--build] [-d] [--with=service]  Start containerized stack")
	fmt.Println("  down [-v]                           Stop and remove containers")
	fmt.Println("  restart [service]                   Restart all or specific service")
	fmt.Println("  status                              Show live container health")
	fmt.Println("  logs [-f] [service]                 Tail service container logs")
	fmt.Println("  pull                                Pull latest container images")
	fmt.Println()
	fmt.Println("⚡ CodeIgniter & Dev Ergonomics:")
	fmt.Println("  spark [command]                     Run CodeIgniter 4 Spark command")
	fmt.Println("  migrate                             Shortcut for 'spark migrate'")
	fmt.Println("  seed [seeder]                       Shortcut for 'spark db:seed'")
	fmt.Println("  composer [command]                  Execute Composer in PHP container")
	fmt.Println("  php [script]                        Execute PHP script in container")
	fmt.Println("  db [options]                        Open interactive MySQL/MariaDB shell")
	fmt.Println("  shell [service]                     Drop into bash/sh shell inside container")
	fmt.Println("  xdebug [on|off|status]              Toggle Xdebug dynamically")
	fmt.Println()
	fmt.Println("🚢 Production & Deployment:")
	fmt.Println("  build [--tag=name] [--no-cache]     Build production multi-stage OCI image")
	fmt.Println("  release / deploy [--skip-health]    Execute safe migration & release pipeline")
	fmt.Println()
	fmt.Println("🛠️ Project Management & QA:")
	fmt.Println("  init [dir] [--force] [--dir=dir]    Bootstrap fresh CI4 AppStarter (Interactive/Custom)")
	fmt.Println("  doctor                              Run host diagnostic checks")
	fmt.Println("  test [options]                      Run PHPUnit test suite")
	fmt.Println("  lint                                Run PHP code style linter")
	fmt.Println("  backup [create|restore]             Fast native backup & restore of src/")
	fmt.Println("  completion [bash|zsh|fish]          Generate shell auto-completions")
	fmt.Println()
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
	appDir := ""

	for _, a := range args {
		if strings.HasPrefix(a, "--version=") {
			tag = strings.TrimPrefix(a, "--version=")
		} else if strings.HasPrefix(a, "--dir=") {
			appDir = strings.TrimPrefix(a, "--dir=")
		} else if strings.HasPrefix(a, "-d=") {
			appDir = strings.TrimPrefix(a, "-d=")
		} else if a == "--force" || a == "-f" {
			force = true
		} else if !strings.HasPrefix(a, "-") && appDir == "" {
			appDir = a
		}
	}

	// Resolve latest version dynamically from GitHub if not specified
	if tag == "" || tag == "latest" {
		fmt.Print("🔍 Checking latest CodeIgniter 4 release from GitHub... ")
		resolvedTag := resolveLatestCI4Release(ctx)
		tag = resolvedTag
		fmt.Printf("(%s)\n", tag)
	}

	pCtx, err := config.FindProjectRoot("")
	var rootDir string
	if err != nil {
		rootDir, _ = os.Getwd()
	} else {
		rootDir = pCtx.RootPath
	}

	// Interactive prompt if appDir not specified via flags/args
	if appDir == "" {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("📂 Nama folder aplikasi CodeIgniter 4 yang diinginkan [default: src]: ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		if input != "" {
			appDir = input
		} else {
			appDir = "src"
		}
	}

	// Sanitize folder name
	appDir = filepath.Clean(appDir)
	targetDir := filepath.Join(rootDir, appDir)

	if !force && dirHasFiles(targetDir) {
		fmt.Printf("⚠️  Folder '%s/' sudah ada dan berisi file. Gunakan --force untuk menimpa.\n", appDir)
		return
	}

	fmt.Printf("📦 Initializing CodeIgniter 4 AppStarter (%s) into '%s/'...\n", tag, appDir)
	tarURL := fmt.Sprintf(AppStarterURL, tag)
	cacheDir := filepath.Join(rootDir, "backups", "cache")
	os.MkdirAll(cacheDir, 0755)
	tarPath := filepath.Join(cacheDir, fmt.Sprintf("appstarter-%s.tar.gz", tag))

	if !fileExists(tarPath) || force {
		fmt.Printf("⬇️  Downloading AppStarter from %s...\n", tarURL)
		if err := downloadFile(tarURL, tarPath); err != nil {
			fmt.Fprintf(os.Stderr, "Download failed: %v\n", err)
			os.Exit(1)
		}
	} else {
		fmt.Printf("⚡ Using cached tarball %s\n", tarPath)
	}

	fmt.Printf("📂 Extracting application skeleton into %s/...\n", appDir)
	if err := extractTarGzStripped(tarPath, targetDir); err != nil {
		fmt.Fprintf(os.Stderr, "Extraction failed: %v\n", err)
		os.Exit(1)
	}

	// Save app dir marker for persistent upward discovery
	_ = os.WriteFile(filepath.Join(rootDir, ".c4ignite-app"), []byte(appDir+"\n"), 0644)

	// Extract embedded templates (docker configs & envs) if missing
	_ = templates.Extract("env", filepath.Join(rootDir, "templates", "env"), false)
	_ = templates.Extract("docker", filepath.Join(rootDir, "docker"), false)

	// Copy dev.env to targetDir/.env if not present
	srcEnv := filepath.Join(targetDir, ".env")
	if !fileExists(srcEnv) {
		devEnv := filepath.Join(rootDir, "templates", "env", "dev.env")
		if fileExists(devEnv) {
			copyFile(devEnv, srcEnv)
			fmt.Printf("📄 Created %s/.env from templates/env/dev.env\n", appDir)
		}
	}

	fmt.Printf("✨ CodeIgniter 4 initialized successfully in '%s/'! Jalankan 'c4ignite up' untuk mulai.\n", appDir)
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

func handleBuild(ctx context.Context, args []string) {
	pCtx, _ := getRunner()
	opts := builder.BuildOptions{
		Tag: "c4ignite-app:latest",
	}

	for _, a := range args {
		if strings.HasPrefix(a, "--tag=") {
			opts.Tag = strings.TrimPrefix(a, "--tag=")
		} else if strings.HasPrefix(a, "-t=") {
			opts.Tag = strings.TrimPrefix(a, "-t=")
		} else if strings.HasPrefix(a, "--target=") {
			opts.Target = strings.TrimPrefix(a, "--target=")
		} else if a == "--no-cache" {
			opts.NoCache = true
		} else if a == "--push" {
			opts.Push = true
		}
	}

	if err := builder.Build(ctx, pCtx, opts); err != nil {
		fmt.Fprintf(os.Stderr, "Build error: %v\n", err)
		os.Exit(1)
	}
}

func handleRelease(ctx context.Context, args []string) {
	pCtx, runner := getRunner()
	opts := release.ReleaseOptions{}

	for _, a := range args {
		if a == "--skip-migration" || a == "--skip-migrate" {
			opts.SkipMigration = true
		} else if a == "--skip-health" {
			opts.SkipHealth = true
		} else if strings.HasPrefix(a, "--health-url=") {
			opts.HealthURL = strings.TrimPrefix(a, "--health-url=")
		}
	}

	if err := release.Execute(ctx, pCtx, runner, opts); err != nil {
		fmt.Fprintf(os.Stderr, "Release error: %v\n", err)
		os.Exit(1)
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
	shell := "bash"
	if len(args) > 0 {
		shell = args[0]
	}
	switch shell {
	case "bash":
		fmt.Println(`_c4ignite() {
    local cur prev opts
    COMPREPLY=()
    cur="${COMP_WORDS[COMP_CWORD]}"
    opts="up down restart status logs pull shell spark composer php db migrate seed test lint xdebug build deploy release backup doctor init version help"
    COMPREPLY=( $(compgen -W "${opts}" -- ${cur}) )
    return 0
}
complete -F _c4ignite c4ignite`)
	case "zsh":
		fmt.Println(`#compdef c4ignite
_c4ignite() {
    local -a commands
    commands=(
        'up:Start containerized stack'
        'down:Stop and remove containers'
        'restart:Restart services'
        'status:Show container health'
        'logs:Tail service logs'
        'spark:Execute CI4 Spark'
        'db:Open interactive MySQL console'
        'migrate:Run database migrations'
        'seed:Seed database'
        'composer:Run Composer in PHP container'
        'php:Run PHP in container'
        'test:Run PHPUnit tests'
        'lint:Run PHP code linter'
        'build:Build production multi-stage OCI container'
        'deploy:Run production release and deployment pipeline'
        'release:Run production release and deployment pipeline'
        'doctor:Run environment diagnostics'
        'backup:Backup or restore src/'
    )
    _describe 'command' commands
}`)
	}
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
