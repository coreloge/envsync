package report

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/user/envsync/internal/envfile"
)

func makeInterpolateResult(substituted, unresolved []string, vars map[string]string) *envfile.InterpolateResult {
	return &envfile.InterpolateResult{
		Substituted: substituted,
		Unresolved:  unresolved,
		Vars:        vars,
	}
}

func TestWriteInterpolation_NoChanges_Text(t *testing.T) {
	res := makeInterpolateResult([]string{}, []string{}, map[string]string{"KEY": "value"})
	var buf bytes.Buffer
	if err := WriteInterpolation(&buf, res, "text"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "No interpolation changes") {
		t.Errorf("expected no-change message, got: %s", buf.String())
	}
}

func TestWriteInterpolation_Substituted_Text(t *testing.T) {
	res := makeInterpolateResult(
		[]string{"DB_URL"},
		[]string{},
		map[string]string{"DB_URL": "postgres://localhost/mydb"},
	)
	var buf bytes.Buffer
	if err := WriteInterpolation(&buf, res, "text"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "DB_URL") {
		t.Errorf("expected DB_URL in output, got: %s", out)
	}
	if !strings.Contains(out, "Substituted") {
		t.Errorf("expected Substituted header, got: %s", out)
	}
}

func TestWriteInterpolation_Unresolved_Text(t *testing.T) {
	res := makeInterpolateResult(
		[]string{},
		[]string{"MISSING_HOST"},
		map[string]string{},
	)
	var buf bytes.Buffer
	if err := WriteInterpolation(&buf, res, "text"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "MISSING_HOST") {
		t.Errorf("expected MISSING_HOST in output, got: %s", out)
	}
	if !strings.Contains(out, "Unresolved") {
		t.Errorf("expected Unresolved header, got: %s", out)
	}
}

func TestWriteInterpolation_JSON_Format(t *testing.T) {
	res := makeInterpolateResult(
		[]string{"URL"},
		[]string{"GHOST"},
		map[string]string{"URL": "https://example.com"},
	)
	var buf bytes.Buffer
	if err := WriteInterpolation(&buf, res, "json"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var out map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if _, ok := out["substituted"]; !ok {
		t.Error("expected 'substituted' key in JSON")
	}
	if _, ok := out["unresolved"]; !ok {
		t.Error("expected 'unresolved' key in JSON")
	}
	if _, ok := out["vars"]; !ok {
		t.Error("expected 'vars' key in JSON")
	}
}
