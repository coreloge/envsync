package report

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"text/tabwriter"

	"github.com/user/envsync/internal/envfile"
)

// WritePromotion writes a human-readable or JSON summary of a PromoteResult.
func WritePromotion(w io.Writer, res envfile.PromoteResult, format string) error {
	switch format {
	case "json":
		return writePromoteJSON(w, res)
	default:
		return writePromoteText(w, res)
	}
}

func writePromoteText(w io.Writer, res envfile.PromoteResult) error {
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	sort.Strings(res.Added)
	sort.Strings(res.Overwritten)
	sort.Strings(res.Skipped)

	if len(res.Added) == 0 && len(res.Overwritten) == 0 && len(res.Skipped) == 0 {
		fmt.Fprintln(tw, "No changes required.")
		return tw.Flush()
	}
	for _, k := range res.Added {
		fmt.Fprintf(tw, "ADDED\t%s\n", k)
	}
	for _, k := range res.Overwritten {
		fmt.Fprintf(tw, "OVERWRITTEN\t%s\n", k)
	}
	for _, k := range res.Skipped {
		fmt.Fprintf(tw, "SKIPPED\t%s\n", k)
	}
	return tw.Flush()
}

func writePromoteJSON(w io.Writer, res envfile.PromoteResult) error {
	payload := map[string][]string{
		"added":       nilToEmpty(res.Added),
		"overwritten": nilToEmpty(res.Overwritten),
		"skipped":     nilToEmpty(res.Skipped),
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(payload)
}

func nilToEmpty(s []string) []string {
	if s == nil {
		return []string{}
	}
	return s
}
