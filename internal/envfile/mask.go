package envfile

import (
	"fmt"
	"strings"
)

// MaskOptions controls how values are masked in output.
type MaskOptions struct {
	// MaskChar is the character used to fill the masked portion (default: "*").
	MaskChar string
	// VisiblePrefix is the number of leading characters to keep visible (default: 0).
	VisiblePrefix int
	// VisibleSuffix is the number of trailing characters to keep visible (default: 0).
	VisibleSuffix int
	// MinMaskLen is the minimum number of mask characters to emit (default: 4).
	MinMaskLen int
}

// DefaultMaskOptions returns sensible defaults for masking.
func DefaultMaskOptions() MaskOptions {
	return MaskOptions{
		MaskChar:      "*",
		VisiblePrefix: 0,
		VisibleSuffix: 0,
		MinMaskLen:    4,
	}
}

// MaskResult holds the output of a Mask operation.
type MaskResult struct {
	// Masked is the new map with sensitive values replaced by mask strings.
	Masked map[string]string
	// MaskedKeys lists the keys whose values were masked.
	MaskedKeys []string
	// SkippedKeys lists the keys that were not sensitive and left unchanged.
	SkippedKeys []string
}

// Mask replaces sensitive values in vars with a masked representation.
// Keys are considered sensitive according to IsSensitiveKey. Non-sensitive
// keys are copied unchanged. The original map is never modified.
func Mask(vars map[string]string, opts MaskOptions) (*MaskResult, error) {
	if vars == nil {
		return nil, fmt.Errorf("mask: vars must not be nil")
	}
	if opts.MaskChar == "" {
		opts.MaskChar = "*"
	}
	if opts.MinMaskLen <= 0 {
		opts.MinMaskLen = 4
	}

	result := &MaskResult{
		Masked:      make(map[string]string, len(vars)),
		MaskedKeys:  []string{},
		SkippedKeys: []string{},
	}

	for k, v := range vars {
		if IsSensitiveKey(k) {
			result.Masked[k] = maskValue(v, opts)
			result.MaskedKeys = append(result.MaskedKeys, k)
		} else {
			result.Masked[k] = v
			result.SkippedKeys = append(result.SkippedKeys, k)
		}
	}

	return result, nil
}

// maskValue applies the masking rules to a single value string.
func maskValue(v string, opts MaskOptions) string {
	l := len(v)
	prefix := opts.VisiblePrefix
	suffix := opts.VisibleSuffix

	if prefix+suffix >= l {
		// Not enough characters to show both ends; mask everything.
		return strings.Repeat(opts.MaskChar, max(opts.MinMaskLen, l))
	}

	maskLen := max(opts.MinMaskLen, l-prefix-suffix)
	return v[:prefix] + strings.Repeat(opts.MaskChar, maskLen) + v[l-suffix:]
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
