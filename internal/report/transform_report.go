package report

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"

	"github.com/user/envsync/internal/envfile"
)

// WriteTransform writes a human-readable or JSON summary of a TransformResult.
func WriteTransform(w io.Writer, result *envfile.TransformResult, format string) error {
	if format == "json" {
		return writeTransformJSON(w, result)
	}
	return writeTransformText(w, result)
}

func writeTransformText(w io.Writer, result *envfile.TransformResult) error {
	changed := sortedTransformSlice(result.Changed)
	unchanged := sortedTransformSlice(result.Unchanged)

	fmt.Fprintf(w, "Transform Summary\n")
	fmt.Fprintf(w, "  Changed:   %d\n", len(changed))
	fmt.Fprintf(w, "  Unchanged: %d\n", len(unchanged))
	fmt.Fprintf(w, "  Total out: %d\n", len(result.Vars))

	if len(changed) > 0 {
		fmt.Fprintf(w, "\nChanged keys (original):\n")
		for _, k := range changed {
			fmt.Fprintf(w, "  ~ %s\n", k)
		}
	}

	if len(unchanged) > 0 {
		fmt.Fprintf(w, "\nUnchanged keys:\n")
		for _, k := range unchanged {
			fmt.Fprintf(w, "  = %s\n", k)
		}
	}

	return nil
}

func writeTransformJSON(w io.Writer, result *envfile.TransformResult) error {
	type payload struct {
		Changed   []string          `json:"changed"`
		Unchanged []string          `json:"unchanged"`
		Vars      map[string]string `json:"vars"`
	}

	p := payload{
		Changed:   sortedTransformSlice(result.Changed),
		Unchanged: sortedTransformSlice(result.Unchanged),
		Vars:      result.Vars,
	}
	if p.Changed == nil {
		p.Changed = []string{}
	}
	if p.Unchanged == nil {
		p.Unchanged = []string{}
	}
	if p.Vars == nil {
		p.Vars = map[string]string{}
	}

	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(p)
}

func sortedTransformSlice(s []string) []string {
	if s == nil {
		return nil
	}
	out := make([]string, len(s))
	copy(out, s)
	sort.Strings(out)
	return out
}
