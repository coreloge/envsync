package envfile

import "fmt"

// PromoteStrategy controls how values are promoted between environments.
type PromoteStrategy int

const (
	// PromoteAddOnly only adds keys missing in the target; never overwrites.
	PromoteAddOnly PromoteStrategy = iota
	// PromoteOverwrite adds missing keys and overwrites mismatched values.
	PromoteOverwrite
)

// PromoteResult describes the outcome of a promotion operation.
type PromoteResult struct {
	Added     []string
	Overwritten []string
	Skipped   []string
}

// Promote copies keys from source into target according to the given strategy.
// It returns a PromoteResult summarising what changed and any error encountered
// when writing the updated target file.
func Promote(sourceVars, targetVars map[string]string, targetPath string, strategy PromoteStrategy) (PromoteResult, error) {
	result := PromoteResult{}
	updated := make(map[string]string, len(targetVars))
	for k, v := range targetVars {
		updated[k] = v
	}

	for key, srcVal := range sourceVars {
		tgtVal, exists := targetVars[key]
		switch {
		case !exists:
			updated[key] = srcVal
			result.Added = append(result.Added, key)
		case exists && tgtVal != srcVal && strategy == PromoteOverwrite:
			updated[key] = srcVal
			result.Overwritten = append(result.Overwritten, key)
		case exists && tgtVal != srcVal && strategy == PromoteAddOnly:
			result.Skipped = append(result.Skipped, key)
		}
	}

	if err := Write(updated, targetPath); err != nil {
		return result, fmt.Errorf("promote: writing target %q: %w", targetPath, err)
	}
	return result, nil
}
