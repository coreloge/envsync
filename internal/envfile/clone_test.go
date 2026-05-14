package envfile

import (
	"sort"
	"testing"
)

func TestClone_AllKeys_NoDst(t *testing.T) {
	src := map[string]string{"A": "1", "B": "2"}
	dst := map[string]string{}
	out, res, err := Clone(src, dst, CloneOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Cloned) != 2 {
		t.Errorf("expected 2 cloned, got %d", len(res.Cloned))
	}
	if len(res.Skipped) != 0 {
		t.Errorf("expected 0 skipped, got %d", len(res.Skipped))
	}
	if out["A"] != "1" || out["B"] != "2" {
		t.Errorf("unexpected output: %v", out)
	}
}

func TestClone_SkipsExistingWithoutOverwrite(t *testing.T) {
	src := map[string]string{"A": "new"}
	dst := map[string]string{"A": "old"}
	out, res, err := Clone(src, dst, CloneOptions{Overwrite: false})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Skipped) != 1 || res.Skipped[0] != "A" {
		t.Errorf("expected A to be skipped, got %v", res.Skipped)
	}
	if out["A"] != "old" {
		t.Errorf("expected old value preserved, got %s", out["A"])
	}
}

func TestClone_OverwriteReplaces(t *testing.T) {
	src := map[string]string{"A": "new"}
	dst := map[string]string{"A": "old"}
	out, res, err := Clone(src, dst, CloneOptions{Overwrite: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Cloned) != 1 {
		t.Errorf("expected 1 cloned, got %d", len(res.Cloned))
	}
	if out["A"] != "new" {
		t.Errorf("expected new value, got %s", out["A"])
	}
}

func TestClone_WithPrefix(t *testing.T) {
	src := map[string]string{"FOO": "bar"}
	dst := map[string]string{}
	out, res, err := Clone(src, dst, CloneOptions{Prefix: "STAGING_"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := out["STAGING_FOO"]; !ok {
		t.Errorf("expected STAGING_FOO in output, got %v", out)
	}
	if len(res.Cloned) != 1 || res.Cloned[0] != "STAGING_FOO" {
		t.Errorf("unexpected cloned keys: %v", res.Cloned)
	}
}

func TestClone_SpecificKeys(t *testing.T) {
	src := map[string]string{"A": "1", "B": "2", "C": "3"}
	dst := map[string]string{}
	out, res, err := Clone(src, dst, CloneOptions{Keys: []string{"A", "C"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	cloned := res.Cloned
	sort.Strings(cloned)
	if len(cloned) != 2 || cloned[0] != "A" || cloned[1] != "C" {
		t.Errorf("expected [A C], got %v", cloned)
	}
	if _, ok := out["B"]; ok {
		t.Errorf("B should not be in output")
	}
}

func TestClone_NilSrcReturnsError(t *testing.T) {
	_, _, err := Clone(nil, map[string]string{}, CloneOptions{})
	if err == nil {
		t.Error("expected error for nil src")
	}
}

func TestClone_NilDstReturnsError(t *testing.T) {
	_, _, err := Clone(map[string]string{}, nil, CloneOptions{})
	if err == nil {
		t.Error("expected error for nil dst")
	}
}
