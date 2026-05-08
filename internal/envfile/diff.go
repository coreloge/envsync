package envfile

// DiffResult holds the result of comparing two env maps.
type DiffResult struct {
	MissingInTarget []string          // keys present in source but not in target
	ExtraInTarget   []string          // keys present in target but not in source
	Mismatched      map[string][2]string // keys present in both but with different values
}

// Diff compares a source env map against a target env map.
// source is treated as the reference (e.g. local or staging).
func Diff(source, target map[string]string) DiffResult {
	result := DiffResult{
		Mismatched: make(map[string][2]string),
	}

	for k, sv := range source {
		tv, ok := target[k]
		if !ok {
			result.MissingInTarget = append(result.MissingInTarget, k)
		} else if sv != tv {
			result.Mismatched[k] = [2]string{sv, tv}
		}
	}

	for k := range target {
		if _, ok := source[k]; !ok {
			result.ExtraInTarget = append(result.ExtraInTarget, k)
		}
	}

	return result
}

// HasDrift returns true if the DiffResult contains any differences.
func (d DiffResult) HasDrift() bool {
	return len(d.MissingInTarget) > 0 || len(d.ExtraInTarget) > 0 || len(d.Mismatched) > 0
}
