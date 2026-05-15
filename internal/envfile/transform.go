package envfile

import (
	"fmt"
	"strings"
)

// TransformFunc is a function that transforms a key-value pair.
// It returns the new key, new value, and whether to include the pair.
type TransformFunc func(key, value string) (newKey, newValue string, keep bool)

// TransformOptions controls how Transform behaves.
type TransformOptions struct {
	// KeyPrefix adds a prefix to all keys.
	KeyPrefix string
	// KeySuffix adds a suffix to all keys.
	KeySuffix string
	// UppercaseKeys converts all keys to uppercase.
	UppercaseKeys bool
	// LowercaseKeys converts all keys to lowercase.
	LowercaseKeys bool
	// ValueTrimSpace trims whitespace from values.
	ValueTrimSpace bool
	// Custom is an optional user-supplied transform applied last.
	Custom TransformFunc
}

// TransformResult holds the outcome of a Transform call.
type TransformResult struct {
	Vars     map[string]string
	Changed  []string
	Unchanged []string
}

// Transform applies a set of transformations to env vars and returns a new map.
func Transform(vars map[string]string, opts TransformOptions) (*TransformResult, error) {
	if vars == nil {
		return nil, fmt.Errorf("transform: vars must not be nil")
	}

	out := make(map[string]string, len(vars))
	result := &TransformResult{
		Vars: out,
	}

	for k, v := range vars {
		newKey := k
		newVal := v

		if opts.UppercaseKeys {
			newKey = strings.ToUpper(newKey)
		} else if opts.LowercaseKeys {
			newKey = strings.ToLower(newKey)
		}

		if opts.KeyPrefix != "" {
			newKey = opts.KeyPrefix + newKey
		}
		if opts.KeySuffix != "" {
			newKey = newKey + opts.KeySuffix
		}
		if opts.ValueTrimSpace {
			newVal = strings.TrimSpace(newVal)
		}

		keep := true
		if opts.Custom != nil {
			newKey, newVal, keep = opts.Custom(newKey, newVal)
		}

		if !keep {
			continue
		}

		out[newKey] = newVal

		if newKey != k || newVal != v {
			result.Changed = append(result.Changed, k)
		} else {
			result.Unchanged = append(result.Unchanged, k)
		}
	}

	return result, nil
}
