package report

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

func TestBuildRedactSummary_CorrectCounts(t *testing.T) {
	original := map[string]string{
		"DB_PASSWORD": "secret",
		"API_KEY":     "key123",
		"APP_NAME":    "myapp",
	}
	redacted := map[string]string{
		"DB_PASSWORD": "***REDACTED***",
		"API_KEY":     "***REDACTED***",
		"APP_NAME":    "myapp",
	}
	s := BuildRedactSummary(original, redacted)

	if s.TotalKeys != 3 {
		t.Errorf("expected 3 total keys, got %d", s.TotalKeys)
	}
	if len(s.RedactedKeys) != 2 {
		t.Errorf("expected 2 redacted keys, got %d", len(s.RedactedKeys))
	}
	if len(s.SafeKeys) != 1 {
		t.Errorf("expected 1 safe key, got %d", len(s.SafeKeys))
	}
}

func TestWriteRedactSummary_TextFormat(t *testing.T) {
	s := RedactSummary{
		TotalKeys:    3,
		RedactedKeys: []string{"DB_PASSWORD", "API_KEY"},
		SafeKeys:     []string{"APP_NAME"},
	}
	var buf bytes.Buffer
	err := WriteRedactSummary(&buf, s, "text")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "Redaction Summary") {
		t.Error("expected header in output")
	}
	if !strings.Contains(out, "DB_PASSWORD") {
		t.Error("expected DB_PASSWORD in redacted list")
	}
	if !strings.Contains(out, "Redacted:") {
		t.Error("expected Redacted: line")
	}
}

func TestWriteRedactSummary_NoRedactedKeys_Text(t *testing.T) {
	s := RedactSummary{TotalKeys: 2, RedactedKeys: nil, SafeKeys: []string{"HOST", "PORT"}}
	var buf bytes.Buffer
	err := WriteRedactSummary(&buf, s, "text")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if strings.Contains(buf.String(), "Redacted keys:") {
		t.Error("should not list redacted keys section when none exist")
	}
}

func TestWriteRedactSummary_JSONFormat(t *testing.T) {
	s := RedactSummary{
		TotalKeys:    2,
		RedactedKeys: []string{"SECRET"},
		SafeKeys:     []string{"HOST"},
	}
	var buf bytes.Buffer
	err := WriteRedactSummary(&buf, s, "json")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var out RedactSummary
	if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if out.TotalKeys != 2 {
		t.Errorf("expected total_keys=2, got %d", out.TotalKeys)
	}
	if len(out.RedactedKeys) != 1 || out.RedactedKeys[0] != "SECRET" {
		t.Errorf("unexpected redacted_keys: %v", out.RedactedKeys)
	}
}
