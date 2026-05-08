package envfile

import (
	"testing"
)

func TestValidate_Clean(t *testing.T) {
	env := map[string]string{
		"APP_ENV":    "production",
		"DB_HOST":    "localhost",
		"SECRET_KEY": "abc123",
	}
	result := Validate(env)
	if result.HasErrors() {
		t.Errorf("expected no errors, got %d: %v", len(result.Errors), result.Errors)
	}
}

func TestValidate_EmptyValue(t *testing.T) {
	env := map[string]string{
		"APP_ENV": "",
		"DB_HOST": "   ",
	}
	result := Validate(env)
	if !result.HasErrors() {
		t.Fatal("expected errors for empty values")
	}
	if len(result.Errors) != 2 {
		t.Errorf("expected 2 errors, got %d", len(result.Errors))
	}
}

func TestValidate_InvalidKeyChars(t *testing.T) {
	env := map[string]string{
		"INVALID-KEY": "value",
		"123STARTS":   "value",
	}
	result := Validate(env)
	if !result.HasErrors() {
		t.Fatal("expected errors for invalid key names")
	}
	if len(result.Errors) != 2 {
		t.Errorf("expected 2 errors, got %d", len(result.Errors))
	}
}

func TestValidate_UnresolvedReference(t *testing.T) {
	env := map[string]string{
		"BASE_URL": "https://${HOST}/api",
	}
	result := Validate(env)
	if !result.HasErrors() {
		t.Fatal("expected error for unresolved variable reference")
	}
	if result.Errors[0].Key != "BASE_URL" {
		t.Errorf("expected error on BASE_URL, got %q", result.Errors[0].Key)
	}
}

func TestValidate_MixedIssues(t *testing.T) {
	env := map[string]string{
		"GOOD_KEY":  "good_value",
		"bad-key":   "value",
		"EMPTY_VAL": "",
	}
	result := Validate(env)
	if len(result.Errors) != 2 {
		t.Errorf("expected 2 errors, got %d: %v", len(result.Errors), result.Errors)
	}
}

func TestValidationError_Error(t *testing.T) {
	e := ValidationError{Key: "MY_KEY", Message: "value is empty"}
	got := e.Error()
	want := `key "MY_KEY": value is empty`
	if got != want {
		t.Errorf("Error() = %q, want %q", got, want)
	}
}
