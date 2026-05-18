package envfile

import (
	"fmt"
	"strings"
)

// FlattenOptions controls how nested key structures are flattened.
type FlattenOptions struct {
	// Separator is the delimiter used to join key segments (default: "_").
	Separator string
	// Prefix is an optional prefix prepended to all output keys.
	Prefix string
	// UppercaseKeys converts all resulting keys to uppercase.
	UppercaseKeys bool
}

// FlattenResult holds the outcome of a Flatten operation.
type FlattenResult struct {
	// Vars is the resulting flattened map.
	Vars map[string]string
	// Renamed maps original keys to their new flattened keys.
	Renamed map[string]string
	// Unchanged lists keys that required no transformation.
	Unchanged []string
}

// Flatten normalises keys that contain a given separator by replacing it with
// a canonical separator, optionally uppercasing keys and adding a prefix.
// It is useful when merging configs from systems that use "." or "/" as
// hierarchy separators into a flat env-var style namespace.
func Flatten(vars map[string]string, inputSep string, opts FlattenOptions) (*FlattenResult, error) {
	if vars == nil {
		return nil, fmt.Errorf("flatten: vars map must not be nil")
	}
	if inputSep == "" {
		return nil, fmt.Errorf("flatten: inputSep must not be empty")
	}

	sep := opts.Separator
	if sep == "" {
		sep = "_"
	}

	result := &FlattenResult{
		Vars:    make(map[string]string, len(vars)),
		Renamed: make(map[string]string),
	}

	for k, v := range vars {
		newKey := k

		if strings.Contains(k, inputSep) {
			newKey = strings.ReplaceAll(k, inputSep, sep)
		}

		if opts.Prefix != "" && !strings.HasPrefix(newKey, opts.Prefix) {
			newKey = opts.Prefix + newKey
		}

		if opts.UppercaseKeys {
			newKey = strings.ToUpper(newKey)
		}

		if newKey != k {
			result.Renamed[k] = newKey
		} else {
			result.Unchanged = append(result.Unchanged, k)
		}

		result.Vars[newKey] = v
	}

	return result, nil
}
