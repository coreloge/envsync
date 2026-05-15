package envfile

import (
	"strings"
	"testing"
)

func TestTransform_NoOptions_ReturnsIdentical(t *testing.T) {
	vars := map[string]string{"FOO": "bar", "BAZ": "qux"}
	res, err := Transform(vars, TransformOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Vars["FOO"] != "bar" || res.Vars["BAZ"] != "qux" {
		t.Errorf("expected unchanged vars, got %v", res.Vars)
	}
	if len(res.Unchanged) != 2 {
		t.Errorf("expected 2 unchanged, got %d", len(res.Unchanged))
	}
}

func TestTransform_UppercaseKeys(t *testing.T) {
	vars := map[string]string{"foo": "1", "bar": "2"}
	res, err := Transform(vars, TransformOptions{UppercaseKeys: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := res.Vars["FOO"]; !ok {
		t.Errorf("expected FOO key, got %v", res.Vars)
	}
	if _, ok := res.Vars["BAR"]; !ok {
		t.Errorf("expected BAR key, got %v", res.Vars)
	}
}

func TestTransform_KeyPrefix(t *testing.T) {
	vars := map[string]string{"HOST": "localhost"}
	res, err := Transform(vars, TransformOptions{KeyPrefix: "APP_"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Vars["APP_HOST"] != "localhost" {
		t.Errorf("expected APP_HOST, got %v", res.Vars)
	}
	if len(res.Changed) != 1 || res.Changed[0] != "HOST" {
		t.Errorf("expected HOST in changed, got %v", res.Changed)
	}
}

func TestTransform_ValueTrimSpace(t *testing.T) {
	vars := map[string]string{"KEY": "  value  "}
	res, err := Transform(vars, TransformOptions{ValueTrimSpace: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Vars["KEY"] != "value" {
		t.Errorf("expected trimmed value, got %q", res.Vars["KEY"])
	}
}

func TestTransform_CustomFunc_FiltersKey(t *testing.T) {
	vars := map[string]string{"KEEP": "yes", "DROP": "no"}
	opts := TransformOptions{
		Custom: func(k, v string) (string, string, bool) {
			return k, v, k != "DROP"
		},
	}
	res, err := Transform(vars, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := res.Vars["DROP"]; ok {
		t.Error("expected DROP to be filtered out")
	}
	if res.Vars["KEEP"] != "yes" {
		t.Error("expected KEEP to remain")
	}
}

func TestTransform_CustomFunc_TransformsValue(t *testing.T) {
	vars := map[string]string{"URL": "http://example.com"}
	opts := TransformOptions{
		Custom: func(k, v string) (string, string, bool) {
			return k, strings.ReplaceAll(v, "http", "https"), true
		},
	}
	res, err := Transform(vars, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Vars["URL"] != "https://example.com" {
		t.Errorf("expected https URL, got %q", res.Vars["URL"])
	}
}

func TestTransform_NilVars_ReturnsError(t *testing.T) {
	_, err := Transform(nil, TransformOptions{})
	if err == nil {
		t.Error("expected error for nil vars")
	}
}
