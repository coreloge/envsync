package envfile

import "fmt"

// DedupStrategy controls how duplicate keys are resolved.
type DedupStrategy int

const (
	// DedupKeepFirst retains the first occurrence of a duplicate key.
	DedupKeepFirst DedupStrategy = iota
	// DedupKeepLast retains the last occurrence of a duplicate key.
	DedupKeepLast
)

// DedupResult holds the outcome of a deduplication pass.
type DedupResult struct {
	// Vars is the deduplicated map of environment variables.
	Vars map[string]string
	// Duplicates lists keys that appeared more than once.
	Duplicates []string
}

// Dedup scans vars for keys that share the same canonical (case-insensitive)
// name and resolves conflicts according to strategy. Because map iteration
// order is non-deterministic, callers should pass an ordered list of keys via
// orderedKeys when the source was a parsed file; if orderedKeys is nil the
// function operates on vars alone and only exact-case duplicates are detected.
func Dedup(vars map[string]string, orderedKeys []string, strategy DedupStrategy) (*DedupResult, error) {
	if vars == nil {
		return nil, fmt.Errorf("dedup: vars must not be nil")
	}

	seen := make(map[string]string)   // canonical lower-key -> first key seen
	dupeSet := make(map[string]bool)  // canonical keys that are duplicated
	result := make(map[string]string)

	keys := orderedKeys
	if len(keys) == 0 {
		for k := range vars {
			keys = append(keys, k)
		}
	}

	for _, k := range keys {
		val, ok := vars[k]
		if !ok {
			continue
		}
		canon := canonicalKey(k)
		if existing, conflict := seen[canon]; conflict {
			dupeSet[canon] = true
			if strategy == DedupKeepLast {
				// Remove the old key entry and write the new one.
				delete(result, existing)
				result[k] = val
				seen[canon] = k
			}
			// DedupKeepFirst: do nothing, first entry already in result.
		} else {
			seen[canon] = k
			result[k] = val
		}
	}

	duplicates := make([]string, 0, len(dupeSet))
	for canon := range dupeSet {
		duplicates = append(duplicates, canon)
	}

	return &DedupResult{
		Vars:       result,
		Duplicates: duplicates,
	}, nil
}

// canonicalKey returns a lowercase version of k for conflict detection.
func canonicalKey(k string) string {
	b := make([]byte, len(k))
	for i := 0; i < len(k); i++ {
		c := k[i]
		if c >= 'A' && c <= 'Z' {
			c += 'a' - 'A'
		}
		b[i] = c
	}
	return string(b)
}
