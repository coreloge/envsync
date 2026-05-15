package report

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/user/envsync/internal/envfile"
)

// WriteCompare writes a human-readable or JSON comparison report.
func WriteCompare(w io.Writer, result envfile.CompareResult, format string) error {
	switch format {
	case "json":
		return writeCompareJSON(w, result)
	default:
		return writeCompareText(w, result)
	}
}

func writeCompareText(w io.Writer, r envfile.CompareResult) error {
	fmt.Fprintf(w, "Comparing %s → %s\n", r.SourceEnv, r.TargetEnv)

	if !r.HasDrift() {
		fmt.Fprintln(w, "✔ No drift detected.")
		return nil
	}

	if len(r.OnlyInSource) > 0 {
		fmt.Fprintf(w, "\nOnly in %s (%d):\n", r.SourceEnv, len(r.OnlyInSource))
		for _, k := range r.OnlyInSource {
			fmt.Fprintf(w, "  - %s\n", k)
		}
	}

	if len(r.OnlyInTarget) > 0 {
		fmt.Fprintf(w, "\nOnly in %s (%d):\n", r.TargetEnv, len(r.OnlyInTarget))
		for _, k := range r.OnlyInTarget {
			fmt.Fprintf(w, "  + %s\n", k)
		}
	}

	if len(r.Mismatched) > 0 {
		fmt.Fprintf(w, "\nMismatched values (%d):\n", len(r.Mismatched))
		for _, k := range r.Mismatched {
			fmt.Fprintf(w, "  ~ %s\n", k)
		}
	}

	fmt.Fprintf(w, "\nMatched: %d key(s)\n", len(r.Matched))
	return nil
}

func writeCompareJSON(w io.Writer, r envfile.CompareResult) error {
	type payload struct {
		SourceEnv    string   `json:"source_env"`
		TargetEnv    string   `json:"target_env"`
		HasDrift     bool     `json:"has_drift"`
		OnlyInSource []string `json:"only_in_source"`
		OnlyInTarget []string `json:"only_in_target"`
		Matched      []string `json:"matched"`
		Mismatched   []string `json:"mismatched"`
	}
	p := payload{
		SourceEnv:    r.SourceEnv,
		TargetEnv:    r.TargetEnv,
		HasDrift:     r.HasDrift(),
		OnlyInSource: nilToEmptyStr(r.OnlyInSource),
		OnlyInTarget: nilToEmptyStr(r.OnlyInTarget),
		Matched:      nilToEmptyStr(r.Matched),
		Mismatched:   nilToEmptyStr(r.Mismatched),
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(p)
}

func nilToEmptyStr(s []string) []string {
	if s == nil {
		return []string{}
	}
	return s
}
