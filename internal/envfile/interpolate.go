package envfile

import (
	"fmt"
	"regexp"
	"strings"
)

// InterpolateOptions controls how variable interpolation behaves.
type InterpolateOptions struct {
	// Strict causes Interpolate to return an error if a referenced variable is not found.
	Strict bool
	// Fallback is used when a referenced variable is missing and Strict is false.
	Fallback string
}

// InterpolateResult holds the output of an interpolation pass.
type InterpolateResult struct {
	// Vars is the new map with all values interpolated.
	Vars map[string]string
	// Substituted lists keys whose values were changed.
	Substituted []string
	// Unresolved lists variable references that could not be resolved.
	Unresolved []string
}

var refPattern = regexp.MustCompile(`\$\{([A-Za-z_][A-Za-z0-9_]*)\}`)

// Interpolate expands ${VAR} references within env values using the same map.
// References can be resolved from the provided vars map itself, enabling
// simple chained substitution (single pass).
func Interpolate(vars map[string]string, opts InterpolateOptions) (*InterpolateResult, error) {
	if vars == nil {
		return nil, fmt.Errorf("interpolate: vars map must not be nil")
	}

	result := &InterpolateResult{
		Vars: make(map[string]string, len(vars)),
		Substituted: []string{},
		Unresolved:  []string{},
	}

	unresolvedSet := map[string]struct{}{}

	for key, value := range vars {
		original := value
		expanded, unresolved := expandValue(value, vars, opts.Fallback)
		result.Vars[key] = expanded

		for _, u := range unresolved {
			if _, seen := unresolvedSet[u]; !seen {
				unresolvedSet[u] = struct{}{}
				result.Unresolved = append(result.Unresolved, u)
			}
		}

		if expanded != original {
			result.Substituted = append(result.Substituted, key)
		}
	}

	if opts.Strict && len(result.Unresolved) > 0 {
		return result, fmt.Errorf("interpolate: unresolved references: %s", strings.Join(result.Unresolved, ", "))
	}

	return result, nil
}

func expandValue(value string, vars map[string]string, fallback string) (string, []string) {
	var unresolved []string
	expanded := refPattern.ReplaceAllStringFunc(value, func(match string) string {
		submatches := refPattern.FindStringSubmatch(match)
		if len(submatches) < 2 {
			return match
		}
		ref := submatches[1]
		if v, ok := vars[ref]; ok {
			return v
		}
		unresolved = append(unresolved, ref)
		return fallback
	})
	return expanded, unresolved
}
