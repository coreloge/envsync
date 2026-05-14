package envfile

import "fmt"

// RenameResult holds the outcome of a rename operation.
type RenameResult struct {
	Renamed  map[string]string // oldKey -> newKey
	Skipped  []string          // keys not found in vars
	Conflict []string          // newKey already exists in vars
}

// RenameOptions controls how rename handles edge cases.
type RenameOptions struct {
	// OverwriteConflict allows renaming even if the target key already exists.
	OverwriteConflict bool
}

// Rename renames keys in vars according to the provided mapping (oldKey -> newKey).
// It returns a RenameResult describing what happened and an updated copy of vars.
func Rename(vars map[string]string, mapping map[string]string, opts RenameOptions) (map[string]string, RenameResult, error) {
	if vars == nil {
		return nil, RenameResult{}, fmt.Errorf("vars must not be nil")
	}
	if mapping == nil {
		return nil, RenameResult{}, fmt.Errorf("mapping must not be nil")
	}

	result := RenameResult{
		Renamed: make(map[string]string),
	}

	out := make(map[string]string, len(vars))
	for k, v := range vars {
		out[k] = v
	}

	for oldKey, newKey := range mapping {
		val, exists := out[oldKey]
		if !exists {
			result.Skipped = append(result.Skipped, oldKey)
			continue
		}
		if _, conflict := out[newKey]; conflict && !opts.OverwriteConflict {
			result.Conflict = append(result.Conflict, oldKey)
			continue
		}
		delete(out, oldKey)
		out[newKey] = val
		result.Renamed[oldKey] = newKey
	}

	return out, result, nil
}
