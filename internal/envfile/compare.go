package envfile

import "sort"

// CompareResult holds the result of comparing two env maps across named environments.
type CompareResult struct {
	SourceEnv    string
	TargetEnv    string
	OnlyInSource []string
	OnlyInTarget []string
	Matched      []string
	Mismatched   []string
}

// CompareOptions controls how comparison is performed.
type CompareOptions struct {
	IgnoreValues bool // if true, only check key presence
}

// Compare performs a detailed comparison between two named env maps.
func Compare(sourceEnv, targetEnv string, source, target map[string]string, opts CompareOptions) CompareResult {
	result := CompareResult{
		SourceEnv: sourceEnv,
		TargetEnv: targetEnv,
	}

	for k, sv := range source {
		if tv, ok := target[k]; !ok {
			result.OnlyInSource = append(result.OnlyInSource, k)
		} else if opts.IgnoreValues || sv == tv {
			result.Matched = append(result.Matched, k)
		} else {
			_ = tv
			result.Mismatched = append(result.Mismatched, k)
		}
	}

	for k := range target {
		if _, ok := source[k]; !ok {
			result.OnlyInTarget = append(result.OnlyInTarget, k)
		}
	}

	sort.Strings(result.OnlyInSource)
	sort.Strings(result.OnlyInTarget)
	sort.Strings(result.Matched)
	sort.Strings(result.Mismatched)

	return result
}

// HasDrift returns true if there are any keys that differ between environments.
func (r CompareResult) HasDrift() bool {
	return len(r.OnlyInSource) > 0 || len(r.OnlyInTarget) > 0 || len(r.Mismatched) > 0
}
