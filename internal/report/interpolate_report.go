package report

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"text/tabwriter"

	"github.com/user/envsync/internal/envfile"
)

// WriteInterpolation writes a human-readable or JSON report of an interpolation result.
func WriteInterpolation(w io.Writer, res *envfile.InterpolateResult, format string) error {
	if format == "json" {
		return writeInterpolateJSON(w, res)
	}
	return writeInterpolateText(w, res)
}

func writeInterpolateText(w io.Writer, res *envfile.InterpolateResult) error {
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)

	sort.Strings(res.Substituted)
	sort.Strings(res.Unresolved)

	if len(res.Substituted) == 0 && len(res.Unresolved) == 0 {
		fmt.Fprintln(tw, "No interpolation changes.")
		return tw.Flush()
	}

	if len(res.Substituted) > 0 {
		fmt.Fprintf(tw, "Substituted (%d):\n", len(res.Substituted))
		for _, k := range res.Substituted {
			fmt.Fprintf(tw, "  ✓\t%s\t→ %s\n", k, res.Vars[k])
		}
	}

	if len(res.Unresolved) > 0 {
		fmt.Fprintf(tw, "Unresolved references (%d):\n", len(res.Unresolved))
		for _, ref := range res.Unresolved {
			fmt.Fprintf(tw, "  ✗\t${%s}\n", ref)
		}
	}

	return tw.Flush()
}

func writeInterpolateJSON(w io.Writer, res *envfile.InterpolateResult) error {
	type payload struct {
		Substituted []string          `json:"substituted"`
		Unresolved  []string          `json:"unresolved"`
		Vars        map[string]string `json:"vars"`
	}

	sub := res.Substituted
	if sub == nil {
		sub = []string{}
	}
	unres := res.Unresolved
	if unres == nil {
		unres = []string{}
	}
	sort.Strings(sub)
	sort.Strings(unres)

	p := payload{
		Substituted: sub,
		Unresolved:  unres,
		Vars:        res.Vars,
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(p)
}
