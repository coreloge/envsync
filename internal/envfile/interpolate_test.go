package envfile

import (
	"strings"
	"testing"
)

func TestInterpolate_NoReferences(t *testing.T) {
	vars := map[string]string{
		"HOST": "localhost",
		"PORT": "5432",
	}
	res, err := Interpolate(vars, InterpolateOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Vars["HOST"] != "localhost" || res.Vars["PORT"] != "5432" {
		t.Errorf("expected unchanged values, got %v", res.Vars)
	}
	if len(res.Substituted) != 0 {
		t.Errorf("expected no substitutions, got %v", res.Substituted)
	}
}

func TestInterpolate_BasicSubstitution(t *testing.T) {
	vars := map[string]string{
		"HOST":    "localhost",
		"DB_HOST": "${HOST}",
	}
	res, err := Interpolate(vars, InterpolateOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Vars["DB_HOST"] != "localhost" {
		t.Errorf("expected DB_HOST=localhost, got %q", res.Vars["DB_HOST"])
	}
	if len(res.Substituted) != 1 || res.Substituted[0] != "DB_HOST" {
		t.Errorf("expected DB_HOST in substituted, got %v", res.Substituted)
	}
}

func TestInterpolate_UnresolvedLenient(t *testing.T) {
	vars := map[string]string{
		"URL": "http://${HOST}:${PORT}/path",
	}
	res, err := Interpolate(vars, InterpolateOptions{Fallback: "MISSING"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Vars["URL"] != "http://MISSING:MISSING/path" {
		t.Errorf("unexpected URL value: %q", res.Vars["URL"])
	}
	if len(res.Unresolved) != 2 {
		t.Errorf("expected 2 unresolved, got %v", res.Unresolved)
	}
}

func TestInterpolate_UnresolvedStrict(t *testing.T) {
	vars := map[string]string{
		"URL": "http://${MISSING_HOST}/api",
	}
	_, err := Interpolate(vars, InterpolateOptions{Strict: true})
	if err == nil {
		t.Fatal("expected error in strict mode, got nil")
	}
	if !strings.Contains(err.Error(), "MISSING_HOST") {
		t.Errorf("expected error to mention MISSING_HOST, got: %v", err)
	}
}

func TestInterpolate_MultipleRefsInValue(t *testing.T) {
	vars := map[string]string{
		"PROTO": "https",
		"HOST":  "example.com",
		"PORT":  "443",
		"URL":   "${PROTO}://${HOST}:${PORT}",
	}
	res, err := Interpolate(vars, InterpolateOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Vars["URL"] != "https://example.com:443" {
		t.Errorf("unexpected URL: %q", res.Vars["URL"])
	}
}

func TestInterpolate_NilVarsReturnsError(t *testing.T) {
	_, err := Interpolate(nil, InterpolateOptions{})
	if err == nil {
		t.Fatal("expected error for nil vars")
	}
}

func TestInterpolate_OriginalUnmodified(t *testing.T) {
	original := map[string]string{
		"BASE": "value",
		"KEY":  "${BASE}_suffix",
	}
	res, _ := Interpolate(original, InterpolateOptions{})
	if original["KEY"] != "${BASE}_suffix" {
		t.Errorf("original map was modified")
	}
	if res.Vars["KEY"] != "value_suffix" {
		t.Errorf("expected interpolated value, got %q", res.Vars["KEY"])
	}
}
