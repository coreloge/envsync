package report

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"

	"github.com/user/envsync/internal/envfile"
)

// WriteDedupReport writes a human-readable or JSON summary of a Dedup operation.
func WriteDedupReport(w io.Writer, result envfile.DedupResult, format string) error {
	switch format {
	case "json":
		return writeDedupJSON(w, result)
	default:
		return writeDedupText(w, result)
	}
}

func writeDedupText(w io.Writer, result envfile.DedupResult) error {
	if len(result.Removed) == 0 {
		_, err := fmt.Fprintln(w, "✔ No duplicate keys found.")
		return err
	}

	_, err := fmt.Fprintf(w, "Removed %d duplicate key(s):\n", len(result.Removed))
	if err != nil {
		return err
	}

	keys := sortedDedupSlice(result.Removed)
	for _, entry := range keys {
		_, err := fmt.Fprintf(w, "  - %-30s (kept: %q, discarded: %q)\n",
			entry.Key, entry.KeptValue, entry.DiscardedValue)
		if err != nil {
			return err
		}
	}

	_, err = fmt.Fprintf(w, "\nRetained %d unique key(s).\n", result.RetainedCount)
	return err
}

type dedupReportJSON struct {
	RemovedCount  int              `json:"removed_count"`
	RetainedCount int              `json:"retained_count"`
	Removed       []dedupEntryJSON `json:"removed"`
}

type dedupEntryJSON struct {
	Key            string `json:"key"`
	KeptValue      string `json:"kept_value"`
	DiscardedValue string `json:"discarded_value"`
}

func writeDedupJSON(w io.Writer, result envfile.DedupResult) error {
	entries := make([]dedupEntryJSON, 0, len(result.Removed))
	for _, e := range sortedDedupSlice(result.Removed) {
		entries = append(entries, dedupEntryJSON{
			Key:            e.Key,
			KeptValue:      e.KeptValue,
			DiscardedValue: e.DiscardedValue,
		})
	}

	payload := dedupReportJSON{
		RemovedCount:  len(result.Removed),
		RetainedCount: result.RetainedCount,
		Removed:       entries,
	}

	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(payload)
}

// sortedDedupSlice returns a stable, alphabetically sorted copy of removed entries.
func sortedDedupSlice(entries []envfile.RemovedEntry) []envfile.RemovedEntry {
	copy_ := make([]envfile.RemovedEntry, len(entries))
	copy(copy_, entries)
	sort.Slice(copy_, func(i, j int) bool {
		return copy_[i].Key < copy_[j].Key
	})
	return copy_
}
