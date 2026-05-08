package envfile

// MergeStrategy defines how conflicts are resolved during a merge.
type MergeStrategy int

const (
	// StrategySourceWins keeps the source value on conflict.
	StrategySourceWins MergeStrategy = iota
	// StrategyTargetWins keeps the target value on conflict.
	StrategyTargetWins
	// StrategyAddMissing only adds keys missing in target; never overwrites.
	StrategyAddMissing
)

// MergeResult holds the outcome of a merge operation.
type MergeResult struct {
	// Merged is the final reconciled map.
	Merged map[string]string
	// Added contains keys that were added to target from source.
	Added []string
	// Overwritten contains keys whose values were overwritten in target.
	Overwritten []string
	// Skipped contains keys that existed in source but were not applied.
	Skipped []string
}

// Merge reconciles source into target using the given strategy.
// It returns a MergeResult describing every action taken.
func Merge(source, target map[string]string, strategy MergeStrategy) MergeResult {
	merged := make(map[string]string, len(target))
	for k, v := range target {
		merged[k] = v
	}

	result := MergeResult{Merged: merged}

	for k, sv := range source {
		tv, exists := target[k]
		switch {
		case !exists:
			// Key is missing in target — always add regardless of strategy.
			merged[k] = sv
			result.Added = append(result.Added, k)
		case sv == tv:
			// Values are identical; nothing to do.
		case strategy == StrategySourceWins:
			merged[k] = sv
			result.Overwritten = append(result.Overwritten, k)
		case strategy == StrategyTargetWins || strategy == StrategyAddMissing:
			result.Skipped = append(result.Skipped, k)
		}
	}

	return result
}
