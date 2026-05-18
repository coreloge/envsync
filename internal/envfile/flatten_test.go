package envfile

import (
	"testing"
)

func TestFlatten_NoNestedKeys(t *testing.T) {
	vars := map[string]string{
		"HOST": "localhost",
		"PORT": "5432",
	}
	res, err := Flatten(vars, ".", FlattenOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Renamed) != 0 {
		t.Errorf("expected no renames, got %v", res.Renamed)
	}
	if len(res.Unchanged) != 2 {
		t.Errorf("expected 2 unchanged, got %d", len(res.Unchanged))
	}
}

func TestFlatten_DotSeparator(t *testing.T) {
	vars := map[string]string{
		"db.host": "localhost",
		"db.port": "5432",
	}
	res, err := Flatten(vars, ".", FlattenOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Vars["db_host"] != "localhost" {
		t.Errorf("expected db_host=localhost, got %q", res.Vars["db_host"])
	}
	if res.Vars["db_port"] != "5432" {
		t.Errorf("expected db_port=5432, got %q", res.Vars["db_port"])
	}
	if res.Renamed["db.host"] != "db_host" {
		t.Errorf("expected rename db.host->db_host, got %q", res.Renamed["db.host"])
	}
}

func TestFlatten_UppercaseKeys(t *testing.T) {
	vars := map[string]string{
		"db.host": "localhost",
	}
	res, err := Flatten(vars, ".", FlattenOptions{UppercaseKeys: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := res.Vars["DB_HOST"]; !ok {
		t.Errorf("expected key DB_HOST, got vars: %v", res.Vars)
	}
}

func TestFlatten_WithPrefix(t *testing.T) {
	vars := map[string]string{
		"host": "localhost",
	}
	res, err := Flatten(vars, ".", FlattenOptions{Prefix: "APP_"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Vars["APP_host"] != "localhost" {
		t.Errorf("expected APP_host=localhost, got %v", res.Vars)
	}
}

func TestFlatten_PrefixNotDuplicated(t *testing.T) {
	vars := map[string]string{
		"APP_host": "localhost",
	}
	res, err := Flatten(vars, ".", FlattenOptions{Prefix: "APP_"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := res.Vars["APP_host"]; !ok {
		t.Errorf("expected APP_host to exist without duplication, got %v", res.Vars)
	}
	if _, dup := res.Vars["APP_APP_host"]; dup {
		t.Errorf("prefix was duplicated unexpectedly")
	}
}

func TestFlatten_NilVarsReturnsError(t *testing.T) {
	_, err := Flatten(nil, ".", FlattenOptions{})
	if err == nil {
		t.Error("expected error for nil vars, got nil")
	}
}

func TestFlatten_EmptyInputSepReturnsError(t *testing.T) {
	_, err := Flatten(map[string]string{"key": "val"}, "", FlattenOptions{})
	if err == nil {
		t.Error("expected error for empty inputSep, got nil")
	}
}

func TestFlatten_CustomOutputSeparator(t *testing.T) {
	vars := map[string]string{
		"db.host": "localhost",
	}
	res, err := Flatten(vars, ".", FlattenOptions{Separator: "__"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Vars["db__host"] != "localhost" {
		t.Errorf("expected db__host=localhost, got %v", res.Vars)
	}
}
