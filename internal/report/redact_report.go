package report

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"text/tabwriter"
)

// RedactSummary describes the result of a redaction pass.
type RedactSummary struct {
	TotalKeys    int      `json:"total_keys"`
	RedactedKeys []string `json:"redacted_keys"`
	SafeKeys     []string `json:"safe_keys"`
}

// WriteRedactSummary writes a human-readable or JSON redaction report to w.
func WriteRedactSummary(w io.Writer, summary RedactSummary, format string) error {
	switch format {
	case "json":
		return writeRedactJSON(w, summary)
	default:
		return writeRedactText(w, summary)
	}
}

func writeRedactText(w io.Writer, s RedactSummary) error {
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintf(tw, "Redaction Summary\n")
	fmt.Fprintf(tw, "-----------------\n")
	fmt.Fprintf(tw, "Total keys:\t%d\n", s.TotalKeys)
	fmt.Fprintf(tw, "Redacted:\t%d\n", len(s.RedactedKeys))
	fmt.Fprintf(tw, "Safe:\t%d\n", len(s.SafeKeys))

	if len(s.RedactedKeys) > 0 {
		sorted := sorted(s.RedactedKeys)
		fmt.Fprintf(tw, "\nRedacted keys:\n")
		for _, k := range sorted {
			fmt.Fprintf(tw, "  - %s\n", k)
		}
	}
	return tw.Flush()
}

func writeRedactJSON(w io.Writer, s RedactSummary) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(s)
}

func sorted(keys []string) []string {
	copy := append([]string{}, keys...)
	sort.Strings(copy)
	return copy
}

// BuildRedactSummary constructs a RedactSummary from original and redacted maps.
func BuildRedactSummary(original, redacted map[string]string) RedactSummary {
	var redactedKeys, safeKeys []string
	for k, origVal := range original {
		if redacted[k] != origVal {
			redactedKeys = append(redactedKeys, k)
		} else {
			safeKeys = append(safeKeys, k)
		}
	}
	return RedactSummary{
		TotalKeys:    len(original),
		RedactedKeys: redactedKeys,
		SafeKeys:     safeKeys,
	}
}
