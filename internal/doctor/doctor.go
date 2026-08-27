package doctor

import (
	"fmt"
	"net"
	"os/exec"
	"time"

	"github.com/neflalabs/c4ignite/internal/config"
)

type CheckResult struct {
	Name    string
	Passed  bool
	Message string
}

// RunDiagnostics executes environmental healthchecks for Docker, ports, and configuration
func RunDiagnostics(pCtx *config.ProjectContext) []CheckResult {
	var results []CheckResult

	// 1. Docker check
	dockerErr := exec.Command("docker", "version").Run()
	if dockerErr == nil {
		results = append(results, CheckResult{
			Name:    "Docker Engine",
			Passed:  true,
			Message: "Docker daemon is running and accessible",
		})
	} else {
		results = append(results, CheckResult{
			Name:    "Docker Engine",
			Passed:  false,
			Message: "Docker daemon is unreachable. Is Docker Desktop or dockerd running?",
		})
	}

	// 2. Docker Compose check
	composeErr := exec.Command("docker", "compose", "version").Run()
	if composeErr != nil {
		composeErr = exec.Command("docker-compose", "version").Run()
	}
	if composeErr == nil {
		results = append(results, CheckResult{
			Name:    "Docker Compose",
			Passed:  true,
			Message: "Docker Compose V2 plugin is active",
		})
	} else {
		results = append(results, CheckResult{
			Name:    "Docker Compose",
			Passed:  false,
			Message: "Docker Compose is missing. Install compose plugin.",
		})
	}

	// 3. Port status checks (8000 for Nginx, 33060 for MariaDB)
	ports := []struct {
		port int
		name string
	}{
		{8000, "HTTP Web (Nginx)"},
		{33060, "Database (MariaDB)"},
	}

	for _, p := range ports {
		addr := fmt.Sprintf("127.0.0.1:%d", p.port)
		conn, err := net.DialTimeout("tcp", addr, 300*time.Millisecond)
		if err != nil {
			results = append(results, CheckResult{
				Name:    fmt.Sprintf("Port %d (%s)", p.port, p.name),
				Passed:  true,
				Message: "Available (Stack is stopped or port is free)",
			})
		} else {
			conn.Close()
			results = append(results, CheckResult{
				Name:    fmt.Sprintf("Port %d (%s)", p.port, p.name),
				Passed:  true,
				Message: "Active (Service is listening and reachable)",
			})
		}
	}

	// 4. Source Directory check
	if pCtx != nil && pCtx.SrcPath != "" {
		results = append(results, CheckResult{
			Name:    "Project Root",
			Passed:  true,
			Message: fmt.Sprintf("Active at %s", pCtx.RootPath),
		})
	}

	return results
}
