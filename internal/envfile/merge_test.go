package envfile

import (
	"sort"
	"testing"
)

func sortedKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func TestMerge_AddMissingKeys(t *testing.T) {
	src := map[string]string{"A": "1", "B": "2"}
	tgt := map[string]string{"A": "1"}

	r := Merge(src, tgt, StrategySourceWins)

	if r.Merged["B"] != "2" {
		t.Errorf("expected B=2, got %q", r.Merged["B"])
	}
	if len(r.Added) != 1 || r.Added[0] != "B" {
		t.Errorf("expected Added=[B], got %v", r.Added)
	}
}

func TestMerge_SourceWins_Overwrite(t *testing.T) {
	src := map[string]string{"KEY": "new"}
	tgt := map[string]string{"KEY": "old"}

	r := Merge(src, tgt, StrategySourceWins)

	if r.Merged["KEY"] != "new" {
		t.Errorf("expected KEY=new, got %q", r.Merged["KEY"])
	}
	if len(r.Overwritten) != 1 || r.Overwritten[0] != "KEY" {
		t.Errorf("expected Overwritten=[KEY], got %v", r.Overwritten)
	}
}

func TestMerge_TargetWins_NoOverwrite(t *testing.T) {
	src := map[string]string{"KEY": "new", "EXTRA": "x"}
	tgt := map[string]string{"KEY": "old"}

	r := Merge(src, tgt, StrategyTargetWins)

	if r.Merged["KEY"] != "old" {
		t.Errorf("expected KEY=old (target wins), got %q", r.Merged["KEY"])
	}
	if len(r.Skipped) != 1 || r.Skipped[0] != "KEY" {
		t.Errorf("expected Skipped=[KEY], got %v", r.Skipped)
	}
	if r.Merged["EXTRA"] != "x" {
		t.Errorf("expected EXTRA to be added, got %q", r.Merged["EXTRA"])
	}
}

func TestMerge_AddMissingStrategy_SkipsConflicts(t *testing.T) {
	src := map[string]string{"A": "src_a", "B": "src_b"}
	tgt := map[string]string{"A": "tgt_a"}

	r := Merge(src, tgt, StrategyAddMissing)

	if r.Merged["A"] != "tgt_a" {
		t.Errorf("expected A=tgt_a, got %q", r.Merged["A"])
	}
	if r.Merged["B"] != "src_b" {
		t.Errorf("expected B=src_b, got %q", r.Merged["B"])
	}
	if len(r.Skipped) != 1 || r.Skipped[0] != "A" {
		t.Errorf("expected Skipped=[A], got %v", r.Skipped)
	}
}

func TestMerge_IdenticalValues_NoAction(t *testing.T) {
	src := map[string]string{"KEY": "same"}
	tgt := map[string]string{"KEY": "same"}

	r := Merge(src, tgt, StrategySourceWins)

	if len(r.Added) != 0 || len(r.Overwritten) != 0 || len(r.Skipped) != 0 {
		t.Errorf("expected no changes for identical values, got %+v", r)
	}
}
