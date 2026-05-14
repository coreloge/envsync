package envfile

import (
	"testing"
)

func TestRename_BasicRename(t *testing.T) {
	vars := map[string]string{"OLD_KEY": "value1", "KEEP": "value2"}
	out, res, err := Rename(vars, map[string]string{"OLD_KEY": "NEW_KEY"}, RenameOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := out["NEW_KEY"]; !ok {
		t.Error("expected NEW_KEY to exist")
	}
	if _, ok := out["OLD_KEY"]; ok {
		t.Error("expected OLD_KEY to be removed")
	}
	if out["KEEP"] != "value2" {
		t.Error("expected KEEP to be unchanged")
	}
	if res.Renamed["OLD_KEY"] != "NEW_KEY" {
		t.Errorf("expected Renamed[OLD_KEY]=NEW_KEY, got %q", res.Renamed["OLD_KEY"])
	}
}

func TestRename_SkipsMissingKey(t *testing.T) {
	vars := map[string]string{"EXISTING": "v"}
	_, res, err := Rename(vars, map[string]string{"MISSING": "NEW"}, RenameOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Skipped) != 1 || res.Skipped[0] != "MISSING" {
		t.Errorf("expected MISSING in Skipped, got %v", res.Skipped)
	}
}

func TestRename_ConflictWithoutOverwrite(t *testing.T) {
	vars := map[string]string{"OLD": "v1", "NEW": "v2"}
	out, res, err := Rename(vars, map[string]string{"OLD": "NEW"}, RenameOptions{OverwriteConflict: false})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Conflict) != 1 || res.Conflict[0] != "OLD" {
		t.Errorf("expected OLD in Conflict, got %v", res.Conflict)
	}
	if out["OLD"] != "v1" {
		t.Error("expected OLD to remain unchanged")
	}
	if out["NEW"] != "v2" {
		t.Error("expected NEW to remain unchanged")
	}
}

func TestRename_ConflictWithOverwrite(t *testing.T) {
	vars := map[string]string{"OLD": "new_val", "NEW": "old_val"}
	out, res, err := Rename(vars, map[string]string{"OLD": "NEW"}, RenameOptions{OverwriteConflict: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Conflict) != 0 {
		t.Errorf("expected no conflicts, got %v", res.Conflict)
	}
	if out["NEW"] != "new_val" {
		t.Errorf("expected NEW=new_val, got %q", out["NEW"])
	}
	if _, ok := out["OLD"]; ok {
		t.Error("expected OLD to be removed")
	}
}

func TestRename_NilVarsReturnsError(t *testing.T) {
	_, _, err := Rename(nil, map[string]string{"A": "B"}, RenameOptions{})
	if err == nil {
		t.Error("expected error for nil vars")
	}
}

func TestRename_OriginalUnmodified(t *testing.T) {
	vars := map[string]string{"OLD": "val"}
	_, _, _ = Rename(vars, map[string]string{"OLD": "NEW"}, RenameOptions{})
	if _, ok := vars["OLD"]; !ok {
		t.Error("expected original vars to be unmodified")
	}
}
