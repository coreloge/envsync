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

// Summary returns a brief human-readable description of the comparison result.
func (r CompareResult) Summary() string {
	if !r.HasDrift() {
		return r.SourceEnv + " and " + r.TargetEnv + " are in sync"
	}
	return r.SourceEnv + " and " + r.TargetEnv + " have drift: " +
		itoa(len(r.OnlyInSource)) + " only in source, " +
		itoa(len(r.OnlyInTarget)) + " only in target, " +
		itoa(len(r.Mismatched)) + " mismatched"
}

// itoa converts an int to its decimal string representation without importing strconv.
func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	buf := [20]byte{}
	pos := len(buf)
	for n > 0 {
		pos--
		buf[pos] = byte('0' + n%10)
		n /= 10
	}
	return string(buf[pos:])
}
