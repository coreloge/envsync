package envfile

import (
	"testing"
)

func TestRedact_SensitiveKeysAreReplaced(t *testing.T) {
	vars := map[string]string{
		"DB_PASSWORD": "supersecret",
		"API_KEY":     "abc123",
		"APP_NAME":    "myapp",
		"AUTH_TOKEN":  "tok_xyz",
	}
	opts := DefaultRedactOptions()
	result := Redact(vars, opts)

	if result["DB_PASSWORD"] != "***REDACTED***" {
		t.Errorf("expected DB_PASSWORD to be redacted, got %q", result["DB_PASSWORD"])
	}
	if result["API_KEY"] != "***REDACTED***" {
		t.Errorf("expected API_KEY to be redacted, got %q", result["API_KEY"])
	}
	if result["AUTH_TOKEN"] != "***REDACTED***" {
		t.Errorf("expected AUTH_TOKEN to be redacted, got %q", result["AUTH_TOKEN"])
	}
	if result["APP_NAME"] != "myapp" {
		t.Errorf("expected APP_NAME to be unchanged, got %q", result["APP_NAME"])
	}
}

func TestRedact_CustomReplacement(t *testing.T) {
	vars := map[string]string{"SECRET_KEY": "s3cr3t", "HOST": "localhost"}
	opts := DefaultRedactOptions()
	opts.Replacement = "<hidden>"
	result := Redact(vars, opts)

	if result["SECRET_KEY"] != "<hidden>" {
		t.Errorf("expected <hidden>, got %q", result["SECRET_KEY"])
	}
	if result["HOST"] != "localhost" {
		t.Errorf("HOST should be unchanged")
	}
}

func TestRedact_OriginalUnmodified(t *testing.T) {
	vars := map[string]string{"PASSWORD": "real"}
	opts := DefaultRedactOptions()
	Redact(vars, opts)
	if vars["PASSWORD"] != "real" {
		t.Error("original map should not be modified")
	}
}

func TestRedact_CaseInsensitiveMatch(t *testing.T) {
	vars := map[string]string{
		"db_password": "secret",
		"MyApiKey":    "key123",
	}
	opts := DefaultRedactOptions()
	result := Redact(vars, opts)

	if result["db_password"] != "***REDACTED***" {
		t.Errorf("expected db_password to be redacted")
	}
	if result["MyApiKey"] != "***REDACTED***" {
		t.Errorf("expected MyApiKey to be redacted")
	}
}

func TestIsSensitiveKey(t *testing.T) {
	opts := DefaultRedactOptions()
	cases := []struct {
		key       string
		expected  bool
	}{
		{"DB_PASSWORD", true},
		{"api_key", true},
		{"SIGNING_KEY", true},
		{"APP_ENV", false},
		{"PORT", false},
	}
	for _, tc := range cases {
		got := IsSensitiveKey(tc.key, opts)
		if got != tc.expected {
			t.Errorf("IsSensitiveKey(%q) = %v, want %v", tc.key, got, tc.expected)
		}
	}
}
