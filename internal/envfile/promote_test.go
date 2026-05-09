package envfile

import (
	"os"
	"path/filepath"
	"sort"
	"testing"
)

func tmpTarget(t *testing.T, vars map[string]string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), ".env.target")
	if err := Write(vars, path); err != nil {
		t.Fatalf("tmpTarget: %v", err)
	}
	return path
}

func TestPromote_AddOnly_AddsMissing(t *testing.T) {
	src := map[string]string{"A": "1", "B": "2", "C": "3"}
	tgt := map[string]string{"A": "1"}
	path := tmpTarget(t, tgt)

	res, err := Promote(src, tgt, path, PromoteAddOnly)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	sort.Strings(res.Added)
	if len(res.Added) != 2 || res.Added[0] != "B" || res.Added[1] != "C" {
		t.Errorf("Added = %v, want [B C]", res.Added)
	}
	if len(res.Overwritten) != 0 {
		t.Errorf("expected no overwrites, got %v", res.Overwritten)
	}
}

func TestPromote_AddOnly_SkipsMismatch(t *testing.T) {
	src := map[string]string{"A": "new"}
	tgt := map[string]string{"A": "old"}
	path := tmpTarget(t, tgt)

	res, err := Promote(src, tgt, path, PromoteAddOnly)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Skipped) != 1 || res.Skipped[0] != "A" {
		t.Errorf("Skipped = %v, want [A]", res.Skipped)
	}
	if len(res.Overwritten) != 0 {
		t.Errorf("expected no overwrites")
	}
}

func TestPromote_Overwrite_UpdatesMismatch(t *testing.T) {
	src := map[string]string{"A": "new", "B": "2"}
	tgt := map[string]string{"A": "old"}
	path := tmpTarget(t, tgt)

	res, err := Promote(src, tgt, path, PromoteOverwrite)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Overwritten) != 1 || res.Overwritten[0] != "A" {
		t.Errorf("Overwritten = %v, want [A]", res.Overwritten)
	}
	if len(res.Added) != 1 || res.Added[0] != "B" {
		t.Errorf("Added = %v, want [B]", res.Added)
	}
}

func TestPromote_WritesFile(t *testing.T) {
	src := map[string]string{"X": "hello"}
	tgt := map[string]string{}
	path := tmpTarget(t, tgt)

	_, err := Promote(src, tgt, path, PromoteAddOnly)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	data, _ := os.ReadFile(path)
	if string(data) != "X=hello\n" {
		t.Errorf("file content = %q, want %q", string(data), "X=hello\n")
	}
}
