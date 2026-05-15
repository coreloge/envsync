package envfile

import "testing"

func TestCompare_NoChanges(t *testing.T) {
	src := map[string]string{"A": "1", "B": "2"}
	dst := map[string]string{"A": "1", "B": "2"}
	r := Compare("local", "staging", src, dst, CompareOptions{})
	if r.HasDrift() {
		t.Fatal("expected no drift")
	}
	if len(r.Matched) != 2 {
		t.Fatalf("expected 2 matched, got %d", len(r.Matched))
	}
}

func TestCompare_OnlyInSource(t *testing.T) {
	src := map[string]string{"A": "1", "B": "2"}
	dst := map[string]string{"A": "1"}
	r := Compare("local", "staging", src, dst, CompareOptions{})
	if !r.HasDrift() {
		t.Fatal("expected drift")
	}
	if len(r.OnlyInSource) != 1 || r.OnlyInSource[0] != "B" {
		t.Fatalf("expected B only in source, got %v", r.OnlyInSource)
	}
}

func TestCompare_OnlyInTarget(t *testing.T) {
	src := map[string]string{"A": "1"}
	dst := map[string]string{"A": "1", "C": "3"}
	r := Compare("local", "staging", src, dst, CompareOptions{})
	if len(r.OnlyInTarget) != 1 || r.OnlyInTarget[0] != "C" {
		t.Fatalf("expected C only in target, got %v", r.OnlyInTarget)
	}
}

func TestCompare_Mismatched(t *testing.T) {
	src := map[string]string{"A": "1", "B": "old"}
	dst := map[string]string{"A": "1", "B": "new"}
	r := Compare("local", "staging", src, dst, CompareOptions{})
	if len(r.Mismatched) != 1 || r.Mismatched[0] != "B" {
		t.Fatalf("expected B mismatched, got %v", r.Mismatched)
	}
}

func TestCompare_IgnoreValues(t *testing.T) {
	src := map[string]string{"A": "1", "B": "old"}
	dst := map[string]string{"A": "1", "B": "new"}
	r := Compare("local", "staging", src, dst, CompareOptions{IgnoreValues: true})
	if r.HasDrift() {
		t.Fatal("expected no drift when ignoring values")
	}
	if len(r.Matched) != 2 {
		t.Fatalf("expected 2 matched, got %d", len(r.Matched))
	}
}

func TestCompare_EnvNamesStored(t *testing.T) {
	r := Compare("local", "prod", map[string]string{}, map[string]string{}, CompareOptions{})
	if r.SourceEnv != "local" || r.TargetEnv != "prod" {
		t.Fatalf("unexpected env names: %q %q", r.SourceEnv, r.TargetEnv)
	}
}
