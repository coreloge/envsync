package report

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/user/envsync/internal/envfile"
)

func TestWriteClone_NoChanges_Text(t *testing.T) {
	var buf bytes.Buffer
	err := WriteClone(&buf, envfile.CloneResult{}, "text")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "No keys cloned") {
		t.Errorf("expected no-op message, got: %s", buf.String())
	}
}

func TestWriteClone_Cloned_Text(t *testing.T) {
	var buf bytes.Buffer
	res := envfile.CloneResult{Cloned: []string{"DB_HOST", "API_KEY"}}
	err := WriteClone(&buf, res, "text")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "Cloned (2)") {
		t.Errorf("expected cloned header, got: %s", out)
	}
	if !strings.Contains(out, "+ API_KEY") {
		t.Errorf("expected API_KEY in output, got: %s", out)
	}
}

func TestWriteClone_Skipped_Text(t *testing.T) {
	var buf bytes.Buffer
	res := envfile.CloneResult{Skipped: []string{"SECRET"}}
	err := WriteClone(&buf, res, "text")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "Skipped (1)") {
		t.Errorf("expected skipped header, got: %s", out)
	}
	if !strings.Contains(out, "already exists") {
		t.Errorf("expected 'already exists' hint, got: %s", out)
	}
}

func TestWriteClone_JSON_Format(t *testing.T) {
	var buf bytes.Buffer
	res := envfile.CloneResult{
		Cloned:  []string{"FOO"},
		Skipped: []string{"BAR"},
	}
	err := WriteClone(&buf, res, "json")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var out map[string][]string
	if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if len(out["cloned"]) != 1 || out["cloned"][0] != "FOO" {
		t.Errorf("unexpected cloned: %v", out["cloned"])
	}
	if len(out["skipped"]) != 1 || out["skipped"][0] != "BAR" {
		t.Errorf("unexpected skipped: %v", out["skipped"])
	}
}
