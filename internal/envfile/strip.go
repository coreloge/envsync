package envfile

import "strings"

// StripOptions controls which keys are removed during a strip operation.
type StripOptions struct {
	// Prefixes removes any key that starts with one of the given prefixes.
	Prefixes []string
	// Keys removes specific keys by exact name.
	Keys []string
	// EmptyValues removes keys whose value is an empty string.
	EmptyValues bool
}

// StripResult holds the outcome of a Strip operation.
type StripResult struct {
	// Stripped contains the keys that were removed.
	Stripped []string
	// Remaining is the resulting map after removal.
	Remaining map[string]string
}

// Strip removes keys from vars according to the provided StripOptions.
// The original map is never modified; a new map is returned inside StripResult.
func Strip(vars map[string]string, opts StripOptions) (*StripResult, error) {
	if vars == nil {
		return nil, fmt.Errorf("strip: vars map must not be nil")
	}

	exactKeys := make(map[string]struct{}, len(opts.Keys))
	for _, k := range opts.Keys {
		exactKeys[k] = struct{}{}
	}

	remaining := make(map[string]string, len(vars))
	var stripped []string

	for k, v := range vars {
		if shouldStrip(k, v, exactKeys, opts) {
			stripped = append(stripped, k)
			continue
		}
		remaining[k] = v
	}

	sort.Strings(stripped)

	return &StripResult{
		Stripped:  stripped,
		Remaining: remaining,
	}, nil
}

func shouldStrip(key, value string, exactKeys map[string]struct{}, opts StripOptions) bool {
	if _, ok := exactKeys[key]; ok {
		return true
	}
	for _, prefix := range opts.Prefixes {
		if strings.HasPrefix(key, prefix) {
			return true
		}
	}
	if opts.EmptyValues && value == "" {
		return true
	}
	return false
}
