package envfile

import (
	"sort"
	"testing"
)

func TestDedup_NoDuplicates(t *testing.T) {
	vars := map[string]string{
		"HOST": "localhost",
		"PORT": "8080",
		"DEBUG": "true",
	}
	res, err := Dedup(vars, nil, DedupKeepFirst)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Duplicates) != 0 {
		t.Errorf("expected no duplicates, got %v", res.Duplicates)
	}
	if len(res.Vars) != 3 {
		t.Errorf("expected 3 vars, got %d", len(res.Vars))
	}
}

func TestDedup_KeepFirst_CaseInsensitive(t *testing.T) {
	vars := map[string]string{
		"DB_HOST": "first",
		"db_host": "second",
	}
	ordered := []string{"DB_HOST", "db_host"}
	res, err := Dedup(vars, ordered, DedupKeepFirst)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Duplicates) != 1 {
		t.Fatalf("expected 1 duplicate, got %v", res.Duplicates)
	}
	if res.Vars["DB_HOST"] != "first" {
		t.Errorf("expected 'first', got %q", res.Vars["DB_HOST"])
	}
	if _, exists := res.Vars["db_host"]; exists {
		t.Error("duplicate key db_host should have been removed")
	}
}

func TestDedup_KeepLast_CaseInsensitive(t *testing.T) {
	vars := map[string]string{
		"API_KEY": "old",
		"api_key": "new",
	}
	ordered := []string{"API_KEY", "api_key"}
	res, err := Dedup(vars, ordered, DedupKeepLast)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Vars["api_key"] != "new" {
		t.Errorf("expected 'new', got %q", res.Vars["api_key"])
	}
	if _, exists := res.Vars["API_KEY"]; exists {
		t.Error("old key API_KEY should have been replaced")
	}
}

func TestDedup_MultipleDuplicates(t *testing.T) {
	vars := map[string]string{
		"FOO": "a",
		"foo": "b",
		"BAR": "x",
		"Bar": "y",
		"BAZ": "z",
	}
	ordered := []string{"FOO", "foo", "BAR", "Bar", "BAZ"}
	res, err := Dedup(vars, ordered, DedupKeepFirst)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	sort.Strings(res.Duplicates)
	if len(res.Duplicates) != 2 {
		t.Errorf("expected 2 duplicates, got %v", res.Duplicates)
	}
	if res.Duplicates[0] != "bar" || res.Duplicates[1] != "foo" {
		t.Errorf("unexpected duplicates: %v", res.Duplicates)
	}
	if len(res.Vars) != 3 {
		t.Errorf("expected 3 unique vars, got %d", len(res.Vars))
	}
}

func TestDedup_NilVarsReturnsError(t *testing.T) {
	_, err := Dedup(nil, nil, DedupKeepFirst)
	if err == nil {
		t.Error("expected error for nil vars")
	}
}
