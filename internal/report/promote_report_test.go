package report

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/user/envsync/internal/envfile"
)

func TestWritePromotion_NoChanges_Text(t *testing.T) {
	res := envfile.PromoteResult{}
	var buf bytes.Buffer
	if err := WritePromotion(&buf, res, "text"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "No changes") {
		t.Errorf("expected 'No changes', got: %s", buf.String())
	}
}

func TestWritePromotion_Added_Text(t *testing.T) {
	res := envfile.PromoteResult{Added: []string{"FOO", "BAR"}}
	var buf bytes.Buffer
	if err := WritePromotion(&buf, res, "text"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "ADDED") || !strings.Contains(out, "FOO") {
		t.Errorf("expected ADDED FOO in output, got: %s", out)
	}
}

func TestWritePromotion_Skipped_Text(t *testing.T) {
	res := envfile.PromoteResult{Skipped: []string{"DB_PASS"}}
	var buf bytes.Buffer
	if err := WritePromotion(&buf, res, "text"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "SKIPPED") {
		t.Errorf("expected SKIPPED in output")
	}
}

func TestWritePromotion_JSON_Format(t *testing.T) {
	res := envfile.PromoteResult{
		Added:       []string{"A"},
		Overwritten: []string{"B"},
		Skipped:     []string{"C"},
	}
	var buf bytes.Buffer
	if err := WritePromotion(&buf, res, "json"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var out map[string][]string
	if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if len(out["added"]) != 1 || out["added"][0] != "A" {
		t.Errorf("added = %v, want [A]", out["added"])
	}
	if len(out["overwritten"]) != 1 || out["overwritten"][0] != "B" {
		t.Errorf("overwritten = %v, want [B]", out["overwritten"])
	}
	if len(out["skipped"]) != 1 || out["skipped"][0] != "C" {
		t.Errorf("skipped = %v, want [C]", out["skipped"])
	}
}
