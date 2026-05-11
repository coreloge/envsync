package envfile

import (
	"regexp"
	"strings"
)

// FilterOptions controls how vars are selected.
type FilterOptions struct {
	// Prefix filters keys that start with the given prefix (case-insensitive).
	Prefix string
	// Pattern filters keys matching the given regex pattern.
	Pattern string
	// Keys filters to an explicit allow-list of key names.
	Keys []string
}

// FilterResult holds the outcome of a filter operation.
type FilterResult struct {
	Matched  map[string]string
	Excluded map[string]string
}

// Filter returns a subset of vars based on the provided FilterOptions.
// If multiple criteria are set, a key must satisfy ALL of them.
func Filter(vars map[string]string, opts FilterOptions) (FilterResult, error) {
	result := FilterResult{
		Matched:  make(map[string]string),
		Excluded: make(map[string]string),
	}

	var re *regexp.Regexp
	if opts.Pattern != "" {
		var err error
		re, err = regexp.Compile(opts.Pattern)
		if err != nil {
			return result, err
		}
	}

	allowSet := make(map[string]struct{}, len(opts.Keys))
	for _, k := range opts.Keys {
		allowSet[k] = struct{}{}
	}

	for k, v := range vars {
		if !matchesFilter(k, opts.Prefix, re, allowSet) {
			result.Excluded[k] = v
			continue
		}
		result.Matched[k] = v
	}

	return result, nil
}

func matchesFilter(key, prefix string, re *regexp.Regexp, allowSet map[string]struct{}) bool {
	if prefix != "" && !strings.HasPrefix(strings.ToUpper(key), strings.ToUpper(prefix)) {
		return false
	}
	if re != nil && !re.MatchString(key) {
		return false
	}
	if len(allowSet) > 0 {
		if _, ok := allowSet[key]; !ok {
			return false
		}
	}
	return true
}
