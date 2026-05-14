package report

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"

	"github.com/user/envsync/internal/envfile"
)

// WriteRename writes a human-readable or JSON report of a rename operation.
func WriteRename(w io.Writer, result envfile.RenameResult, format string) error {
	switch format {
	case "json":
		return writeRenameJSON(w, result)
	default:
		return writeRenameText(w, result)
	}
}

func writeRenameText(w io.Writer, result envfile.RenameResult) error {
	if len(result.Renamed) == 0 && len(result.Skipped) == 0 && len(result.Conflict) == 0 {
		_, err := fmt.Fprintln(w, "No keys renamed.")
		return err
	}

	if len(result.Renamed) > 0 {
		fmt.Fprintln(w, "Renamed:")
		keys := make([]string, 0, len(result.Renamed))
		for k := range result.Renamed {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, old := range keys {
			fmt.Fprintf(w, "  %s -> %s\n", old, result.Renamed[old])
		}
	}

	if len(result.Skipped) > 0 {
		sort.Strings(result.Skipped)
		fmt.Fprintln(w, "Skipped (not found):")
		for _, k := range result.Skipped {
			fmt.Fprintf(w, "  %s\n", k)
		}
	}

	if len(result.Conflict) > 0 {
		sort.Strings(result.Conflict)
		fmt.Fprintln(w, "Conflicts (target key already exists):")
		for _, k := range result.Conflict {
			fmt.Fprintf(w, "  %s\n", k)
		}
	}
	return nil
}

func writeRenameJSON(w io.Writer, result envfile.RenameResult) error {
	type payload struct {
		Renamed  map[string]string `json:"renamed"`
		Skipped  []string          `json:"skipped"`
		Conflict []string          `json:"conflict"`
	}
	p := payload{
		Renamed:  result.Renamed,
		Skipped:  nilToEmptyStrSlice(result.Skipped),
		Conflict: nilToEmptyStrSlice(result.Conflict),
	}
	if p.Renamed == nil {
		p.Renamed = map[string]string{}
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(p)
}

func nilToEmptyStrSlice(s []string) []string {
	if s == nil {
		return []string{}
	}
	return s
}
