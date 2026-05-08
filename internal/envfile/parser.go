package envfile

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// EnvMap represents a set of environment variable key-value pairs.
type EnvMap map[string]string

// Parse reads a .env file from the given path and returns an EnvMap.
// Lines beginning with '#' are treated as comments and skipped.
// Empty lines are also skipped.
func Parse(path string) (EnvMap, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("envfile: open %q: %w", path, err)
	}
	defer f.Close()

	env := make(EnvMap)
	scanner := bufio.NewScanner(f)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())

		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		key, value, err := parseLine(line)
		if err != nil {
			return nil, fmt.Errorf("envfile: %q line %d: %w", path, lineNum, err)
		}

		env[key] = value
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("envfile: scanning %q: %w", path, err)
	}

	return env, nil
}

// parseLine splits a single KEY=VALUE line into its components.
func parseLine(line string) (key, value string, err error) {
	parts := strings.SplitN(line, "=", 2)
	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid format %q: expected KEY=VALUE", line)
	}

	key = strings.TrimSpace(parts[0])
	if key == "" {
		return "", "", fmt.Errorf("empty key in line %q", line)
	}

	value = strings.TrimSpace(parts[1])
	// Strip optional surrounding quotes.
	if len(value) >= 2 {
		if (value[0] == '"' && value[len(value)-1] == '"') ||
			(value[0] == '\'' && value[len(value)-1] == '\'') {
			value = value[1 : len(value)-1]
		}
	}

	return key, value, nil
}
