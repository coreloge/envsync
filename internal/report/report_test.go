package report

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/envsync/internal/envfile"
)

func TestWriter_NoDrift_Text(t *testing.T) {
	var buf bytes.Buffer
	w := NewWriter(&buf, FormatText)
	result := envfile.DiffResult{Mismatched: map[string][2]string{}}
	if err := w.Write("local", "staging", result); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "No drift") {
		t.Errorf("expected no-drift message, got: %s", buf.String())
	}
}

func TestWriter_MissingKeys_Text(t *testing.T) {
	var buf bytes.Buffer
	w := NewWriter(&buf, FormatText)
	result := envfile.DiffResult{
		MissingInTarget: []string{"SECRET_KEY"},
		Mismatched:      map[string][2]string{},
	}
	w.Write("local", "staging", result)
	output := buf.String()
	if !strings.Contains(output, "SECRET_KEY") {
		t.Errorf("expected SECRET_KEY in output, got: %s", output)
	}
	if !strings.Contains(output, "Missing in staging") {
		t.Errorf("expected missing section header, got: %s", output)
	}
}

func TestWriter_Mismatched_Text(t *testing.T) {
	var buf bytes.Buffer
	w := NewWriter(&buf, FormatText)
	result := envfile.DiffResult{
		Mismatched: map[string][2]string{
			"DB_HOST": {"localhost", "db.prod.example.com"},
		},
	}
	w.Write("local", "production", result)
	output := buf.String()
	if !strings.Contains(output, "DB_HOST") {
		t.Errorf("expected DB_HOST in output, got: %s", output)
	}
	if !strings.Contains(output, "localhost") {
		t.Errorf("expected source value in output")
	}
}

func TestWriter_JSON_Format(t *testing.T) {
	var buf bytes.Buffer
	w := NewWriter(&buf, FormatJSON)
	result := envfile.DiffResult{
		MissingInTarget: []string{"FOO"},
		Mismatched:      map[string][2]string{},
	}
	w.Write("local", "staging", result)
	output := buf.String()
	if !strings.Contains(output, `"source":"local"`) {
		t.Errorf("expected JSON source field, got: %s", output)
	}
	if !strings.Contains(output, `"FOO"`) {
		t.Errorf("expected FOO in JSON output, got: %s", output)
	}
}
