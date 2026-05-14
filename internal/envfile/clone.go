package envfile

import "fmt"

// CloneOptions controls how variables are cloned between environments.
type CloneOptions struct {
	// Keys limits cloning to specific keys; empty means all keys.
	Keys []string
	// Overwrite replaces existing keys in the destination.
	Overwrite bool
	// Prefix is prepended to every cloned key name.
	Prefix string
}

// CloneResult holds the outcome of a Clone operation.
type CloneResult struct {
	Cloned  []string
	Skipped []string
}

// Clone copies variables from src into dst according to opts.
// It returns a CloneResult describing which keys were cloned or skipped.
func Clone(src, dst map[string]string, opts CloneOptions) (map[string]string, CloneResult, error) {
	if src == nil {
		return nil, CloneResult{}, fmt.Errorf("clone: source vars must not be nil")
	}
	if dst == nil {
		return nil, CloneResult{}, fmt.Errorf("clone: destination vars must not be nil")
	}

	result := CloneResult{}
	output := make(map[string]string, len(dst))
	for k, v := range dst {
		output[k] = v
	}

	sourceKeys := resolveKeys(src, opts.Keys)

	for _, key := range sourceKeys {
		val, ok := src[key]
		if !ok {
			continue
		}
		destKey := opts.Prefix + key
		if _, exists := output[destKey]; exists && !opts.Overwrite {
			result.Skipped = append(result.Skipped, destKey)
			continue
		}
		output[destKey] = val
		result.Cloned = append(result.Cloned, destKey)
	}

	return output, result, nil
}

// resolveKeys returns the keys to operate on: opts.Keys if provided, else all src keys.
func resolveKeys(src map[string]string, keys []string) []string {
	if len(keys) > 0 {
		return keys
	}
	out := make([]string, 0, len(src))
	for k := range src {
		out = append(out, k)
	}
	return out
}
