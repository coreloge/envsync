package envfile

import "sort"

// DiffResult holds the outcome of comparing two EnvMaps.
type DiffResult struct {
	// MissingInTarget are keys present in source but absent in target.
	MissingInTarget []string
	// ExtraInTarget are keys present in target but absent in source.
	ExtraInTarget []string
	// ValueMismatch are keys present in both but with differing values.
	ValueMismatch []KeyDiff
}

// KeyDiff describes a value difference for a single key.
type KeyDiff struct {
	Key         string
	SourceValue string
	TargetValue string
}

// IsClean returns true when there are no differences between the two maps.
func (d DiffResult) IsClean() bool {
	return len(d.MissingInTarget) == 0 &&
		len(d.ExtraInTarget) == 0 &&
		len(d.ValueMismatch) == 0
}

// Diff compares a source EnvMap against a target EnvMap and returns a DiffResult.
// It does not compare values for keys listed in ignoreKeys.
func Diff(source, target EnvMap, ignoreKeys ...string) DiffResult {
	ignore := make(map[string]struct{}, len(ignoreKeys))
	for _, k := range ignoreKeys {
		ignore[k] = struct{}{}
	}

	var result DiffResult

	for k, sv := range source {
		if _, skip := ignore[k]; skip {
			continue
		}
		tv, exists := target[k]
		if !exists {
			result.MissingInTarget = append(result.MissingInTarget, k)
		} else if sv != tv {
			result.ValueMismatch = append(result.ValueMismatch, KeyDiff{
				Key:         k,
				SourceValue: sv,
				TargetValue: tv,
			})
		}
	}

	for k := range target {
		if _, skip := ignore[k]; skip {
			continue
		}
		if _, exists := source[k]; !exists {
			result.ExtraInTarget = append(result.ExtraInTarget, k)
		}
	}

	sort.Strings(result.MissingInTarget)
	sort.Strings(result.ExtraInTarget)
	sort.Slice(result.ValueMismatch, func(i, j int) bool {
		return result.ValueMismatch[i].Key < result.ValueMismatch[j].Key
	})

	return result
}
