package report

import (
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/user/envsync/internal/envfile"
)

// Format controls the output style of the report.
type Format string

const (
	FormatText Format = "text"
	FormatJSON Format = "json"
)

// Writer renders a DiffResult to an io.Writer.
type Writer struct {
	out    io.Writer
	format Format
}

// NewWriter creates a new report Writer.
func NewWriter(out io.Writer, format Format) *Writer {
	return &Writer{out: out, format: format}
}

// Write renders the diff result with the source and target labels.
func (w *Writer) Write(source, target string, result envfile.DiffResult) error {
	if w.format == FormatJSON {
		return w.writeJSON(source, target, result)
	}
	return w.writeText(source, target, result)
}

func (w *Writer) writeText(source, target string, result envfile.DiffResult) error {
	if !result.HasDrift() {
		_, err := fmt.Fprintf(w.out, "✓ No drift detected between %s and %s\n", source, target)
		return err
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Drift report: %s → %s\n", source, target))
	sb.WriteString(strings.Repeat("-", 40) + "\n")

	if len(result.MissingInTarget) > 0 {
		sort.Strings(result.MissingInTarget)
		sb.WriteString(fmt.Sprintf("Missing in %s (%d):\n", target, len(result.MissingInTarget)))
		for _, k := range result.MissingInTarget {
			sb.WriteString(fmt.Sprintf("  - %s\n", k))
		}
	}

	if len(result.ExtraInTarget) > 0 {
		sort.Strings(result.ExtraInTarget)
		sb.WriteString(fmt.Sprintf("Extra in %s (%d):\n", target, len(result.ExtraInTarget)))
		for _, k := range result.ExtraInTarget {
			sb.WriteString(fmt.Sprintf("  + %s\n", k))
		}
	}

	if len(result.Mismatched) > 0 {
		keys := make([]string, 0, len(result.Mismatched))
		for k := range result.Mismatched {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		sb.WriteString(fmt.Sprintf("Mismatched values (%d):\n", len(result.Mismatched)))
		for _, k := range keys {
			pair := result.Mismatched[k]
			sb.WriteString(fmt.Sprintf("  ~ %s: %q → %q\n", k, pair[0], pair[1]))
		}
	}

	_, err := fmt.Fprint(w.out, sb.String())
	return err
}

func (w *Writer) writeJSON(source, target string, result envfile.DiffResult) error {
	_, err := fmt.Fprintf(w.out,
		`{"source":%q,"target":%q,"missing_in_target":%s,"extra_in_target":%s,"mismatched":%s}\n`,
		source, target,
		jsonStringSlice(result.MissingInTarget),
		jsonStringSlice(result.ExtraInTarget),
		jsonMismatched(result.Mismatched),
	)
	return err
}

func jsonStringSlice(ss []string) string {
	if len(ss) == 0 {
		return "[]"
	}
	var sb strings.Builder
	sb.WriteString("[")
	for i, s := range ss {
		if i > 0 {
			sb.WriteString(",")
		}
		sb.WriteString(fmt.Sprintf("%q", s))
	}
	sb.WriteString("]")
	return sb.String()
}

func jsonMismatched(m map[string][2]string) string {
	if len(m) == 0 {
		return "{}"
	}
	var sb strings.Builder
	sb.WriteString("{")
	i := 0
	for k, pair := range m {
		if i > 0 {
			sb.WriteString(",")
		}
		sb.WriteString(fmt.Sprintf("%q:{\"source\":%q,\"target\":%q}", k, pair[0], pair[1]))
		i++
	}
	sb.WriteString("}")
	return sb.String()
}
