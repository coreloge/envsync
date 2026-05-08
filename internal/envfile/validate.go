package envfile

import (
	"fmt"
	"regexp"
	"strings"
)

// ValidationError represents a single validation issue found in an env map.
type ValidationError struct {
	Key     string
	Message string
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("key %q: %s", e.Key, e.Message)
}

// ValidationResult holds all errors found during validation.
type ValidationResult struct {
	Errors []ValidationError
}

func (r *ValidationResult) HasErrors() bool {
	return len(r.Errors) > 0
}

func (r *ValidationResult) Add(key, message string) {
	r.Errors = append(r.Errors, ValidationError{Key: key, Message: message})
}

var validKeyPattern = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*$`)

// Validate checks an env map for common issues:
// - empty keys or values
// - keys with invalid characters
// - values with unresolved variable references (${...})
func Validate(env map[string]string) *ValidationResult {
	result := &ValidationResult{}

	for key, value := range env {
		if key == "" {
			result.Add(key, "key must not be empty")
			continue
		}

		if !validKeyPattern.MatchString(key) {
			result.Add(key, fmt.Sprintf("key contains invalid characters (must match %s)", validKeyPattern))
		}

		if strings.TrimSpace(value) == "" {
			result.Add(key, "value is empty or whitespace-only")
		}

		if strings.Contains(value, "${") {
			result.Add(key, "value contains unresolved variable reference")
		}
	}

	return result
}
