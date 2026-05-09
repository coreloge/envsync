package envfile

import (
	"regexp"
	"strings"
)

// RedactOptions controls how sensitive values are redacted.
type RedactOptions struct {
	// Patterns is a list of key name patterns (case-insensitive) that trigger redaction.
	Patterns []string
	// Replacement is the string used to replace sensitive values.
	Replacement string
}

// DefaultRedactOptions returns sensible defaults for secret detection.
func DefaultRedactOptions() RedactOptions {
	return RedactOptions{
		Patterns: []string{
			"password", "passwd", "secret", "token",
			"api_key", "apikey", "private_key", "auth",
			"credential", "access_key", "signing_key",
		},
		Replacement: "***REDACTED***",
	}
}

// Redact returns a copy of vars with sensitive values replaced.
// Keys are matched case-insensitively against each pattern as a substring.
func Redact(vars map[string]string, opts RedactOptions) map[string]string {
	compiled := compilePatterns(opts.Patterns)
	replacement := opts.Replacement
	if replacement == "" {
		replacement = "***REDACTED***"
	}

	result := make(map[string]string, len(vars))
	for k, v := range vars {
		if isSensitive(k, compiled) {
			result[k] = replacement
		} else {
			result[k] = v
		}
	}
	return result
}

// IsSensitiveKey reports whether a key name matches any sensitive pattern.
func IsSensitiveKey(key string, opts RedactOptions) bool {
	return isSensitive(key, compilePatterns(opts.Patterns))
}

func isSensitive(key string, patterns []*regexp.Regexp) bool {
	lower := strings.ToLower(key)
	for _, p := range patterns {
		if p.MatchString(lower) {
			return true
		}
	}
	return false
}

func compilePatterns(patterns []string) []*regexp.Regexp {
	result := make([]*regexp.Regexp, 0, len(patterns))
	for _, p := range patterns {
		result = append(result, regexp.MustCompile(regexp.QuoteMeta(strings.ToLower(p))))
	}
	return result
}
