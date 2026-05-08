package envfile

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestNewSnapshot_CopiesVars(t *testing.T) {
	orig := map[string]string{"FOO": "bar", "BAZ": "qux"}
	s := NewSnapshot("test", orig)

	if s.Label != "test" {
		t.Errorf("expected label 'test', got %q", s.Label)
	}
	if len(s.Vars) != 2 {
		t.Errorf("expected 2 vars, got %d", len(s.Vars))
	}

	// Mutating original should not affect snapshot
	orig["FOO"] = "mutated"
	if s.Vars["FOO"] != "bar" {
		t.Errorf("snapshot vars should be independent of source map")
	}
}

func TestNewSnapshot_TimestampSet(t *testing.T) {
	before := time.Now().UTC()
	s := NewSnapshot("ts-test", map[string]string{})
	after := time.Now().UTC()

	if s.Timestamp.Before(before) || s.Timestamp.After(after) {
		t.Errorf("snapshot timestamp %v out of expected range [%v, %v]", s.Timestamp, before, after)
	}
}

func TestSaveAndLoadSnapshot_RoundTrip(t *testing.T) {
	vars := map[string]string{
		"APP_ENV":  "staging",
		"DB_HOST":  "localhost",
		"LOG_LEVEL": "debug",
	}
	s := NewSnapshot("staging", vars)

	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "snapshot.json")

	if err := SaveSnapshot(path, s); err != nil {
		t.Fatalf("SaveSnapshot failed: %v", err)
	}

	loaded, err := LoadSnapshot(path)
	if err != nil {
		t.Fatalf("LoadSnapshot failed: %v", err)
	}

	if loaded.Label != s.Label {
		t.Errorf("label mismatch: got %q, want %q", loaded.Label, s.Label)
	}
	for k, v := range vars {
		if loaded.Vars[k] != v {
			t.Errorf("var %q: got %q, want %q", k, loaded.Vars[k], v)
		}
	}
}

func TestLoadSnapshot_FileNotFound(t *testing.T) {
	_, err := LoadSnapshot("/nonexistent/path/snapshot.json")
	if err == nil {
		t.Error("expected error for missing file, got nil")
	}
}

func TestSaveSnapshot_InvalidPath(t *testing.T) {
	s := NewSnapshot("bad", map[string]string{})
	err := SaveSnapshot("/nonexistent/dir/snap.json", s)
	if err == nil {
		t.Error("expected error for invalid path, got nil")
	}
}

func TestDiffSnapshot_DetectsDrift(t *testing.T) {
	base := NewSnapshot("local", map[string]string{"FOO": "1", "BAR": "2"})
	target := NewSnapshot("staging", map[string]string{"FOO": "1", "BAZ": "3"})

	result := DiffSnapshot(base, target)

	if !result.HasDrift {
		t.Error("expected drift to be detected")
	}
	if len(result.MissingInTarget) != 1 || result.MissingInTarget[0] != "BAR" {
		t.Errorf("expected BAR missing in target, got %v", result.MissingInTarget)
	}
	if len(result.ExtraInTarget) != 1 || result.ExtraInTarget[0] != "BAZ" {
		t.Errorf("expected BAZ extra in target, got %v", result.ExtraInTarget)
	}
}

func init() {
	// ensure os is used
	_ = os.Stderr
}
