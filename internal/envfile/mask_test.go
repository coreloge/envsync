package envfile

import (
	"testing"
)

func TestMask_SensitiveKeysAreMasked(t *testing.T) {
	vars := map[string]string{
		"API_SECRET": "supersecret",
		"APP_ENV":    "production",
	}
	opts := DefaultMaskOptions()
	res, err := Mask(vars, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Masked["API_SECRET"] == "supersecret" {
		t.Error("expected API_SECRET to be masked")
	}
	if res.Masked["APP_ENV"] != "production" {
		t.Errorf("expected APP_ENV unchanged, got %q", res.Masked["APP_ENV"])
	}
}

func TestMask_VisiblePrefixAndSuffix(t *testing.T) {
	vars := map[string]string{"DB_PASSWORD": "abcdefghij"}
	opts := MaskOptions{
		MaskChar:      "*",
		VisiblePrefix: 2,
		VisibleSuffix: 2,
		MinMaskLen:    4,
	}
	res, err := Mask(vars, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := res.Masked["DB_PASSWORD"]
	if got[:2] != "ab" {
		t.Errorf("expected prefix 'ab', got %q", got[:2])
	}
	if got[len(got)-2:] != "ij" {
		t.Errorf("expected suffix 'ij', got %q", got[len(got)-2:])
	}
}

func TestMask_MinMaskLen(t *testing.T) {
	vars := map[string]string{"SECRET_KEY": "ab"}
	opts := MaskOptions{MaskChar: "*", MinMaskLen: 6}
	res, err := Mask(vars, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := res.Masked["SECRET_KEY"]
	if len(got) < 6 {
		t.Errorf("expected at least 6 mask chars, got %q (len %d)", got, len(got))
	}
}

func TestMask_NilVarsReturnsError(t *testing.T) {
	_, err := Mask(nil, DefaultMaskOptions())
	if err == nil {
		t.Error("expected error for nil vars")
	}
}

func TestMask_OriginalUnmodified(t *testing.T) {
	vars := map[string]string{"TOKEN": "mytoken123"}
	opts := DefaultMaskOptions()
	_, err := Mask(vars, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if vars["TOKEN"] != "mytoken123" {
		t.Error("original map was modified")
	}
}

func TestMask_MaskedAndSkippedKeysCounted(t *testing.T) {
	vars := map[string]string{
		"SECRET":  "s",
		"API_KEY": "k",
		"HOST":    "localhost",
		"PORT":    "8080",
	}
	res, err := Mask(vars, DefaultMaskOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.MaskedKeys) != 2 {
		t.Errorf("expected 2 masked keys, got %d", len(res.MaskedKeys))
	}
	if len(res.SkippedKeys) != 2 {
		t.Errorf("expected 2 skipped keys, got %d", len(res.SkippedKeys))
	}
}
