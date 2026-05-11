package envfile

import (
	"testing"
)

var sampleVars = map[string]string{
	"DB_HOST":     "localhost",
	"DB_PORT":     "5432",
	"APP_NAME":    "envsync",
	"APP_VERSION": "1.0.0",
	"SECRET_KEY":  "abc123",
}

func TestFilter_NoOptions_ReturnsAll(t *testing.T) {
	result, err := Filter(sampleVars, FilterOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Matched) != len(sampleVars) {
		t.Errorf("expected %d matched, got %d", len(sampleVars), len(result.Matched))
	}
	if len(result.Excluded) != 0 {
		t.Errorf("expected 0 excluded, got %d", len(result.Excluded))
	}
}

func TestFilter_ByPrefix(t *testing.T) {
	result, err := Filter(sampleVars, FilterOptions{Prefix: "DB_"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Matched) != 2 {
		t.Errorf("expected 2 matched, got %d", len(result.Matched))
	}
	if _, ok := result.Matched["DB_HOST"]; !ok {
		t.Error("expected DB_HOST in matched")
	}
	if _, ok := result.Matched["DB_PORT"]; !ok {
		t.Error("expected DB_PORT in matched")
	}
}

func TestFilter_ByPattern(t *testing.T) {
	result, err := Filter(sampleVars, FilterOptions{Pattern: "^APP_"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Matched) != 2 {
		t.Errorf("expected 2 matched, got %d", len(result.Matched))
	}
}

func TestFilter_ByKeys(t *testing.T) {
	result, err := Filter(sampleVars, FilterOptions{Keys: []string{"SECRET_KEY", "APP_NAME"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Matched) != 2 {
		t.Errorf("expected 2 matched, got %d", len(result.Matched))
	}
	if _, ok := result.Matched["SECRET_KEY"]; !ok {
		t.Error("expected SECRET_KEY in matched")
	}
}

func TestFilter_InvalidPattern_ReturnsError(t *testing.T) {
	_, err := Filter(sampleVars, FilterOptions{Pattern: "[invalid"})
	if err == nil {
		t.Error("expected error for invalid regex pattern")
	}
}

func TestFilter_CombinedPrefixAndKeys(t *testing.T) {
	// Only DB_ prefix AND in allow-list
	result, err := Filter(sampleVars, FilterOptions{
		Prefix: "DB_",
		Keys:   []string{"DB_HOST", "APP_NAME"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Only DB_HOST satisfies both prefix and allow-list
	if len(result.Matched) != 1 {
		t.Errorf("expected 1 matched, got %d", len(result.Matched))
	}
	if _, ok := result.Matched["DB_HOST"]; !ok {
		t.Error("expected DB_HOST in matched")
	}
}
