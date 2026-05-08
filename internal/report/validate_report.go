package report

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/user/envsync/internal/envfile"
)

// WriteValidation writes a human-readable or JSON report of validation errors.
func (w *Writer) WriteValidation(result *envfile.ValidationResult, label string) error {
	if w.format == "json" {
		return w.writeValidationJSON(result, label)
	}
	return w.writeValidationText(result, label)
}

func (w *Writer) writeValidationText(result *envfile.ValidationResult, label string) error {
	if !result.HasErrors() {
		_, err := fmt.Fprintf(w.out, "[%s] validation passed — no issues found\n", label)
		return err
	}

	_, err := fmt.Fprintf(w.out, "[%s] validation failed — %d issue(s) found:\n", label, len(result.Errors))
	if err != nil {
		return err
	}

	for _, e := range result.Errors {
		_, err = fmt.Fprintf(w.out, "  - %s\n", e.Error())
		if err != nil {
			return err
		}
	}
	return nil
}

type validationJSONOutput struct {
	Label  string   `json:"label"`
	Passed bool     `json:"passed"`
	Errors []string `json:"errors,omitempty"`
}

func (w *Writer) writeValidationJSON(result *envfile.ValidationResult, label string) error {
	out := validationJSONOutput{
		Label:  label,
		Passed: !result.HasErrors(),
	}
	for _, e := range result.Errors {
		out.Errors = append(out.Errors, e.Error())
	}
	return writeJSON(w.out, out)
}

func writeJSON(w io.Writer, v any) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(v)
}
