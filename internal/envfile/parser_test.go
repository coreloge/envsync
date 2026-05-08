package envfile

import (
	"os"
	"path/filepath"
	"testing"
)

func writeTempEnv(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("failed to write temp env file: %v", err)
	}
	return path
}

func TestParse_BasicKeyValue(t *testing.T) {
	path := writeTempEnv(t, "APP_ENV=production\nPORT=8080\n")
	env, err := Parse(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if env["APP_ENV"] != "production" {
		t.Errorf("APP_ENV: got %q, want %q", env["APP_ENV"], "production")
	}
	if env["PORT"] != "8080" {
		t.Errorf("PORT: got %q, want %q", env["PORT"], "8080")
	}
}

func TestParse_CommentsAndBlankLines(t *testing.T) {
	content := "# this is a comment\n\nDB_HOST=localhost\n"
	path := writeTempEnv(t, content)
	env, err := Parse(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(env) != 1 {
		t.Errorf("expected 1 key, got %d", len(env))
	}
	if env["DB_HOST"] != "localhost" {
		t.Errorf("DB_HOST: got %q, want %q", env["DB_HOST"], "localhost")
	}
}

func TestParse_QuotedValues(t *testing.T) {
	path := writeTempEnv(t, `SECRET="my secret value"\nTOKEN='abc123'\n`)
	env, err := Parse(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if env["SECRET"] != "my secret value" {
		t.Errorf("SECRET: got %q, want %q", env["SECRET"], "my secret value")
	}
	if env["TOKEN"] != "abc123" {
		t.Errorf("TOKEN: got %q, want %q", env["TOKEN"], "abc123")
	}
}

func TestParse_InvalidLine(t *testing.T) {
	path := writeTempEnv(t, "INVALID_LINE_NO_EQUALS\n")
	_, err := Parse(path)
	if err == nil {
		t.Fatal("expected error for invalid line, got nil")
	}
}

func TestParse_FileNotFound(t *testing.T) {
	_, err := Parse("/nonexistent/path/.env")
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}
