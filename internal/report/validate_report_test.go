package report

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/envsync/internal/envfile"
)

func TestWriteValidation_Passed_Text(t *testing.T) {
	var buf bytes.Buffer
	w := NewWriter(&buf, "text")

	result := &envfile.ValidationResult{}
	if err := w.WriteValidation(result, "local"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got := buf.String()
	if !strings.Contains(got, "validation passed") {
		t.Errorf("expected 'validation passed' in output, got: %q", got)
	}
	if !strings.Contains(got, "[local]") {
		t.Errorf("expected label '[local]' in output, got: %q", got)
	}
}

func TestWriteValidation_Errors_Text(t *testing.T) {
	var buf bytes.Buffer
	w := NewWriter(&buf, "text")

	result := &envfile.ValidationResult{}
	result.Add("BAD-KEY", "key contains invalid characters")
	result.Add("EMPTY_VAL", "value is empty or whitespace-only")

	if err := w.WriteValidation(result, "staging"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got := buf.String()
	if !strings.Contains(got, "validation failed") {
		t.Errorf("expected 'validation failed' in output, got: %q", got)
	}
	if !strings.Contains(got, "2 issue(s)") {
		t.Errorf("expected issue count in output, got: %q", got)
	}
	if !strings.Contains(got, "BAD-KEY") {
		t.Errorf("expected BAD-KEY in output, got: %q", got)
	}
}

func TestWriteValidation_JSON_Passed(t *testing.T) {
	var buf bytes.Buffer
	w := NewWriter(&buf, "json")

	result := &envfile.ValidationResult{}
	if err := w.WriteValidation(result, "prod"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got := buf.String()
	if !strings.Contains(got, `"passed": true`) {
		t.Errorf("expected passed:true in JSON output, got: %q", got)
	}
	if !strings.Contains(got, `"label": "prod"`) {
		t.Errorf("expected label in JSON output, got: %q", got)
	}
}

func TestWriteValidation_JSON_WithErrors(t *testing.T) {
	var buf bytes.Buffer
	w := NewWriter(&buf, "json")

	result := &envfile.ValidationResult{}
	result.Add("MY_VAR", "value contains unresolved variable reference")

	if err := w.WriteValidation(result, "prod"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got := buf.String()
	if !strings.Contains(got, `"passed": false`) {
		t.Errorf("expected passed:false in JSON output, got: %q", got)
	}
	if !strings.Contains(got, `"errors"`) {
		t.Errorf("expected errors field in JSON output, got: %q", got)
	}
}
