package envfile

import (
	"sort"
	"strings"
)

// SortOrder defines the ordering strategy for environment variable keys.
type SortOrder int

const (
	// SortAscending sorts keys A→Z (default).
	SortAscending SortOrder = iota
	// SortDescending sorts keys Z→A.
	SortDescending
	// SortByLength sorts keys shortest→longest.
	SortByLength
)

// SortOptions controls how Sort behaves.
type SortOptions struct {
	// Order specifies the sort direction or strategy.
	Order SortOrder
	// CaseInsensitive normalises keys to lowercase for comparison.
	CaseInsensitive bool
}

// SortResult holds the outcome of a Sort operation.
type SortResult struct {
	// Sorted is the new map with keys in the requested order (map itself is
	// unordered; SortedKeys carries the canonical order).
	Sorted map[string]string
	// SortedKeys is the slice of keys in the requested order.
	SortedKeys []string
	// OriginalKeys is the slice of keys in the original insertion order
	// (alphabetical as returned by the parser).
	OriginalKeys []string
}

// Sort returns a SortResult containing the vars map reordered according to
// opts. The original map is never modified.
func Sort(vars map[string]string, opts SortOptions) (*SortResult, error) {
	if vars == nil {
		return nil, fmt.Errorf("sort: vars map must not be nil")
	}

	keys := make([]string, 0, len(vars))
	for k := range vars {
		keys = append(keys, k)
	}

	// Capture original order (stable baseline: alphabetical).
	original := make([]string, len(keys))
	copy(original, keys)
	sort.Strings(original)

	switch opts.Order {
	case SortDescending:
		if opts.CaseInsensitive {
			sort.Slice(keys, func(i, j int) bool {
				return strings.ToLower(keys[i]) > strings.ToLower(keys[j])
			})
		} else {
			sort.Sort(sort.Reverse(sort.StringSlice(keys)))
		}
	case SortByLength:
		sort.Slice(keys, func(i, j int) bool {
			if len(keys[i]) != len(keys[j]) {
				return len(keys[i]) < len(keys[j])
			}
			return keys[i] < keys[j]
		})
	default: // SortAscending
		if opts.CaseInsensitive {
			sort.Slice(keys, func(i, j int) bool {
				return strings.ToLower(keys[i]) < strings.ToLower(keys[j])
			})
		} else {
			sort.Strings(keys)
		}
	}

	sorted := make(map[string]string, len(vars))
	for k, v := range vars {
		sorted[k] = v
	}

	return &SortResult{
		Sorted:       sorted,
		SortedKeys:   keys,
		OriginalKeys: original,
	}, nil
}
