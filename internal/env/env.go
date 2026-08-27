package env

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Map represents key-value pairs of environment variables
type Map map[string]string

// Load parses an env file into key-value pairs
func Load(path string) (Map, error) {
	result := make(Map)
	file, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return result, nil
		}
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			k := strings.TrimSpace(parts[0])
			v := strings.TrimSpace(parts[1])
			v = strings.Trim(v, `"'`)
			result[k] = v
		}
	}
	return result, scanner.Err()
}

// SetKey updates or appends a key in an env file while preserving existing content and comments
func SetKey(path, key, value string) error {
	var lines []string
	found := false

	file, err := os.Open(path)
	if err == nil {
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()
			trimmed := strings.TrimSpace(line)
			if !strings.HasPrefix(trimmed, "#") && strings.Contains(trimmed, "=") {
				parts := strings.SplitN(trimmed, "=", 2)
				if strings.TrimSpace(parts[0]) == key {
					lines = append(lines, fmt.Sprintf("%s=%s", key, value))
					found = true
					continue
				}
			}
			lines = append(lines, line)
		}
		file.Close()
	}

	if !found {
		lines = append(lines, fmt.Sprintf("%s=%s", key, value))
	}

	return os.WriteFile(path, []byte(strings.Join(lines, "\n")+"\n"), 0644)
}
