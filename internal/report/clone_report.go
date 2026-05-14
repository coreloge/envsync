package report

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"

	"github.com/user/envsync/internal/envfile"
)

// WriteClone writes a human-readable or JSON report for a clone operation.
func WriteClone(w io.Writer, result envfile.CloneResult, format string) error {
	switch format {
	case "json":
		return writeCloneJSON(w, result)
	default:
		return writeCloneText(w, result)
	}
}

func writeCloneText(w io.Writer, result envfile.CloneResult) error {
	cloned := sortedCopySlice(result.Cloned)
	skipped := sortedCopySlice(result.Skipped)

	if len(cloned) == 0 && len(skipped) == 0 {
		_, err := fmt.Fprintln(w, "No keys cloned.")
		return err
	}

	if len(cloned) > 0 {
		fmt.Fprintf(w, "Cloned (%d):\n", len(cloned))
		for _, k := range cloned {
			fmt.Fprintf(w, "  + %s\n", k)
		}
	}

	if len(skipped) > 0 {
		fmt.Fprintf(w, "Skipped (%d):\n", len(skipped))
		for _, k := range skipped {
			fmt.Fprintf(w, "  ~ %s (already exists)\n", k)
		}
	}
	return nil
}

type cloneJSONReport struct {
	Cloned  []string `json:"cloned"`
	Skipped []string `json:"skipped"`
}

func writeCloneJSON(w io.Writer, result envfile.CloneResult) error {
	cloned := result.Cloned
	if cloned == nil {
		cloned = []string{}
	}
	skipped := result.Skipped
	if skipped == nil {
		skipped = []string{}
	}
	sort.Strings(cloned)
	sort.Strings(skipped)
	rep := cloneJSONReport{Cloned: cloned, Skipped: skipped}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(rep)
}

func sortedCopySlice(s []string) []string {
	out := make([]string, len(s))
	copy(out, s)
	sort.Strings(out)
	return out
}
