package envfile

import (
	"testing"
)

func TestStrip_ByExactKey(t *testing.T) {
	vars := map[string]string{
		"DB_HOST":     "localhost",
		"DB_PASSWORD": "secret",
		"APP_PORT":    "8080",
	}
	res, err := Strip(vars, StripOptions{Keys: []string{"DB_PASSWORD"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Stripped) != 1 || res.Stripped[0] != "DB_PASSWORD" {
		t.Errorf("expected DB_PASSWORD stripped, got %v", res.Stripped)
	}
	if _, ok := res.Remaining["DB_PASSWORD"]; ok {
		t.Error("DB_PASSWORD should not be in remaining")
	}
	if res.Remaining["DB_HOST"] != "localhost" {
		t.Error("DB_HOST should remain")
	}
}

func TestStrip_ByPrefix(t *testing.T) {
	vars := map[string]string{
		"DB_HOST":  "localhost",
		"DB_PORT":  "5432",
		"APP_PORT": "8080",
	}
	res, err := Strip(vars, StripOptions{Prefixes: []string{"DB_"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Stripped) != 2 {
		t.Errorf("expected 2 stripped, got %d: %v", len(res.Stripped), res.Stripped)
	}
	if _, ok := res.Remaining["APP_PORT"]; !ok {
		t.Error("APP_PORT should remain")
	}
}

func TestStrip_EmptyValues(t *testing.T) {
	vars := map[string]string{
		"KEY_A": "value",
		"KEY_B": "",
		"KEY_C": "",
	}
	res, err := Strip(vars, StripOptions{EmptyValues: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Stripped) != 2 {
		t.Errorf("expected 2 stripped, got %d", len(res.Stripped))
	}
	if _, ok := res.Remaining["KEY_A"]; !ok {
		t.Error("KEY_A should remain")
	}
}

func TestStrip_NilVarsReturnsError(t *testing.T) {
	_, err := Strip(nil, StripOptions{})
	if err == nil {
		t.Error("expected error for nil vars")
	}
}

func TestStrip_OriginalUnmodified(t *testing.T) {
	vars := map[string]string{
		"KEY_A": "a",
		"KEY_B": "b",
	}
	_, err := Strip(vars, StripOptions{Keys: []string{"KEY_A"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := vars["KEY_A"]; !ok {
		t.Error("original map should not be modified")
	}
}

func TestStrip_NoOptions_ReturnsAll(t *testing.T) {
	vars := map[string]string{"A": "1", "B": "2"}
	res, err := Strip(vars, StripOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Stripped) != 0 {
		t.Errorf("expected nothing stripped, got %v", res.Stripped)
	}
	if len(res.Remaining) != 2 {
		t.Errorf("expected 2 remaining, got %d", len(res.Remaining))
	}
}
