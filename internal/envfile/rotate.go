package envfile

import (
	"fmt"
	"time"
)

// RotateOptions controls how key rotation is performed.
type RotateOptions struct {
	// Keys is the explicit list of keys to rotate. If empty, all keys are rotated.
	Keys []string
	// BackupSuffix is appended to the backup file path. Defaults to ".bak".
	BackupSuffix string
}

// RotateResult holds the outcome of a rotation operation.
type RotateResult struct {
	Rotated   []string
	Skipped   []string
	Timestamp time.Time
}

// Rotate replaces the values of the specified keys (or all keys) with the
// provided replacements map. Keys present in vars but absent from replacements
// are skipped. Returns a RotateResult describing what changed.
func Rotate(vars map[string]string, replacements map[string]string, opts RotateOptions) (map[string]string, RotateResult, error) {
	if vars == nil {
		return nil, RotateResult{}, fmt.Errorf("vars must not be nil")
	}

	targetKeys := opts.Keys
	if len(targetKeys) == 0 {
		for k := range vars {
			targetKeys = append(targetKeys, k)
		}
	}

	result := make(map[string]string, len(vars))
	for k, v := range vars {
		result[k] = v
	}

	var rotated, skipped []string
	for _, key := range targetKeys {
		if _, exists := vars[key]; !exists {
			skipped = append(skipped, key)
			continue
		}
		newVal, hasReplacement := replacements[key]
		if !hasReplacement {
			skipped = append(skipped, key)
			continue
		}
		result[key] = newVal
		rotated = append(rotated, key)
	}

	return result, RotateResult{
		Rotated:   rotated,
		Skipped:   skipped,
		Timestamp: time.Now().UTC(),
	}, nil
}
