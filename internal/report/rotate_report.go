package report

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"time"

	"github.com/user/envsync/internal/envfile"
)

type rotateJSON struct {
	Timestamp string   `json:"timestamp"`
	Rotated   []string `json:"rotated"`
	Skipped   []string `json:"skipped"`
	Total     int      `json:"total_rotated"`
}

// WriteRotation writes a human-readable or JSON summary of a rotation result.
func WriteRotation(w io.Writer, res envfile.RotateResult, format string) error {
	if format == "json" {
		return writeRotateJSON(w, res)
	}
	return writeRotateText(w, res)
}

func writeRotateText(w io.Writer, res envfile.RotateResult) error {
	rotated := sortedCopy(res.Rotated)
	skipped := sortedCopy(res.Skipped)

	fmt.Fprintf(w, "Rotation completed at %s\n", res.Timestamp.Format(time.RFC3339))
	fmt.Fprintf(w, "  Rotated : %d key(s)\n", len(rotated))
	for _, k := range rotated {
		fmt.Fprintf(w, "    + %s\n", k)
	}
	if len(skipped) > 0 {
		fmt.Fprintf(w, "  Skipped : %d key(s)\n", len(skipped))
		for _, k := range skipped {
			fmt.Fprintf(w, "    - %s\n", k)
		}
	}
	return nil
}

func writeRotateJSON(w io.Writer, res envfile.RotateResult) error {
	payload := rotateJSON{
		Timestamp: res.Timestamp.Format(time.RFC3339),
		Rotated:   sortedCopy(res.Rotated),
		Skipped:   sortedCopy(res.Skipped),
		Total:     len(res.Rotated),
	}
	if payload.Rotated == nil {
		payload.Rotated = []string{}
	}
	if payload.Skipped == nil {
		payload.Skipped = []string{}
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(payload)
}

func sortedCopy(in []string) []string {
	out := make([]string, len(in))
	copy(out, in)
	sort.Strings(out)
	return out
}
