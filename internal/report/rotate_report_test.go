package report

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/user/envsync/internal/envfile"
)

func makeRotateResult(rotated, skipped []string) envfile.RotateResult {
	return envfile.RotateResult{
		Rotated:   rotated,
		Skipped:   skipped,
		Timestamp: time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC),
	}
}

func TestWriteRotation_Text_RotatedAndSkipped(t *testing.T) {
	res := makeRotateResult([]string{"DB_PASSWORD", "API_KEY"}, []string{"APP_ENV"})
	var buf bytes.Buffer
	if err := WriteRotation(&buf, res, "text"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "Rotated") {
		t.Error("expected 'Rotated' in output")
	}
	if !strings.Contains(out, "DB_PASSWORD") {
		t.Error("expected DB_PASSWORD in output")
	}
	if !strings.Contains(out, "Skipped") {
		t.Error("expected 'Skipped' in output")
	}
	if !strings.Contains(out, "APP_ENV") {
		t.Error("expected APP_ENV in skipped")
	}
}

func TestWriteRotation_Text_NoSkipped(t *testing.T) {
	res := makeRotateResult([]string{"TOKEN"}, nil)
	var buf bytes.Buffer
	if err := WriteRotation(&buf, res, "text"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if strings.Contains(out, "Skipped") {
		t.Error("should not mention Skipped when none")
	}
}

func TestWriteRotation_JSON_Format(t *testing.T) {
	res := makeRotateResult([]string{"SECRET"}, []string{"UNRELATED"})
	var buf bytes.Buffer
	if err := WriteRotation(&buf, res, "json"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var payload map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &payload); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if payload["total_rotated"].(float64) != 1 {
		t.Errorf("expected total_rotated=1")
	}
	rotated := payload["rotated"].([]interface{})
	if len(rotated) != 1 || rotated[0].(string) != "SECRET" {
		t.Errorf("unexpected rotated list: %v", rotated)
	}
}

func TestWriteRotation_JSON_EmptySlices(t *testing.T) {
	res := makeRotateResult(nil, nil)
	var buf bytes.Buffer
	if err := WriteRotation(&buf, res, "json"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var payload map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &payload); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if payload["rotated"] == nil || payload["skipped"] == nil {
		t.Error("rotated and skipped should be empty arrays, not null")
	}
}
