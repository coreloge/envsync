package report

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/user/envsync/internal/envfile"
)

func makeCompareResult(src, tgt string, onlySrc, onlyTgt, matched, mismatched []string) envfile.CompareResult {
	return envfile.CompareResult{
		SourceEnv:    src,
		TargetEnv:    tgt,
		OnlyInSource: onlySrc,
		OnlyInTarget: onlyTgt,
		Matched:      matched,
		Mismatched:   mismatched,
	}
}

func TestWriteCompare_NoDrift_Text(t *testing.T) {
	r := makeCompareResult("local", "staging", nil, nil, []string{"A", "B"}, nil)
	var buf bytes.Buffer
	if err := WriteCompare(&buf, r, "text"); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(buf.String(), "No drift") {
		t.Errorf("expected no drift message, got: %s", buf.String())
	}
}

func TestWriteCompare_OnlyInSource_Text(t *testing.T) {
	r := makeCompareResult("local", "staging", []string{"SECRET"}, nil, []string{"A"}, nil)
	var buf bytes.Buffer
	if err := WriteCompare(&buf, r, "text"); err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	if !strings.Contains(out, "SECRET") {
		t.Errorf("expected SECRET in output, got: %s", out)
	}
	if !strings.Contains(out, "Only in local") {
		t.Errorf("expected 'Only in local', got: %s", out)
	}
}

func TestWriteCompare_Mismatched_Text(t *testing.T) {
	r := makeCompareResult("local", "prod", nil, nil, nil, []string{"DB_HOST"})
	var buf bytes.Buffer
	if err := WriteCompare(&buf, r, "text"); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(buf.String(), "DB_HOST") {
		t.Errorf("expected DB_HOST in mismatched output")
	}
}

func TestWriteCompare_JSON_Format(t *testing.T) {
	r := makeCompareResult("local", "prod", []string{"EXTRA"}, nil, []string{"A"}, []string{"B"})
	var buf bytes.Buffer
	if err := WriteCompare(&buf, r, "json"); err != nil {
		t.Fatal(err)
	}
	var out map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if out["source_env"] != "local" {
		t.Errorf("expected source_env=local, got %v", out["source_env"])
	}
	if out["has_drift"] != true {
		t.Errorf("expected has_drift=true")
	}
}
