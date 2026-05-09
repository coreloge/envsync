package envfile

import (
	"testing"
)

func TestRotate_AllKeys(t *testing.T) {
	vars := map[string]string{
		"DB_PASSWORD": "old_pass",
		"API_KEY":     "old_key",
	}
	replacements := map[string]string{
		"DB_PASSWORD": "new_pass",
		"API_KEY":     "new_key",
	}
	out, res, err := Rotate(vars, replacements, RotateOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["DB_PASSWORD"] != "new_pass" {
		t.Errorf("expected new_pass, got %s", out["DB_PASSWORD"])
	}
	if out["API_KEY"] != "new_key" {
		t.Errorf("expected new_key, got %s", out["API_KEY"])
	}
	if len(res.Rotated) != 2 {
		t.Errorf("expected 2 rotated, got %d", len(res.Rotated))
	}
	if len(res.Skipped) != 0 {
		t.Errorf("expected 0 skipped, got %d", len(res.Skipped))
	}
}

func TestRotate_SpecificKeys(t *testing.T) {
	vars := map[string]string{
		"DB_PASSWORD": "old_pass",
		"API_KEY":     "old_key",
		"APP_ENV":     "production",
	}
	replacements := map[string]string{
		"DB_PASSWORD": "new_pass",
		"API_KEY":     "new_key",
	}
	out, res, err := Rotate(vars, replacements, RotateOptions{Keys: []string{"DB_PASSWORD"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["DB_PASSWORD"] != "new_pass" {
		t.Errorf("expected new_pass, got %s", out["DB_PASSWORD"])
	}
	if out["API_KEY"] != "old_key" {
		t.Errorf("API_KEY should be unchanged")
	}
	if len(res.Rotated) != 1 || res.Rotated[0] != "DB_PASSWORD" {
		t.Errorf("expected only DB_PASSWORD rotated")
	}
}

func TestRotate_MissingReplacementSkips(t *testing.T) {
	vars := map[string]string{"SECRET": "old"}
	out, res, err := Rotate(vars, map[string]string{}, RotateOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["SECRET"] != "old" {
		t.Errorf("expected old value preserved")
	}
	if len(res.Skipped) != 1 {
		t.Errorf("expected 1 skipped, got %d", len(res.Skipped))
	}
}

func TestRotate_NilVarsReturnsError(t *testing.T) {
	_, _, err := Rotate(nil, map[string]string{}, RotateOptions{})
	if err == nil {
		t.Error("expected error for nil vars")
	}
}

func TestRotate_OriginalUnmodified(t *testing.T) {
	vars := map[string]string{"TOKEN": "original"}
	replacements := map[string]string{"TOKEN": "rotated"}
	_, _, _ = Rotate(vars, replacements, RotateOptions{})
	if vars["TOKEN"] != "original" {
		t.Error("original map should not be modified")
	}
}
