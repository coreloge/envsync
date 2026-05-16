package envfile

import (
	"testing"
)

func TestSort_Ascending(t *testing.T) {
	vars := map[string]string{
		"ZEBRA": "1",
		"ALPHA": "2",
		"MANGO": "3",
	}
	res, err := Sort(vars, SortOptions{Order: SortAscending})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := []string{"ALPHA", "MANGO", "ZEBRA"}
	for i, k := range res.SortedKeys {
		if k != want[i] {
			t.Errorf("SortedKeys[%d] = %q, want %q", i, k, want[i])
		}
	}
}

func TestSort_Descending(t *testing.T) {
	vars := map[string]string{
		"ZEBRA": "1",
		"ALPHA": "2",
		"MANGO": "3",
	}
	res, err := Sort(vars, SortOptions{Order: SortDescending})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := []string{"ZEBRA", "MANGO", "ALPHA"}
	for i, k := range res.SortedKeys {
		if k != want[i] {
			t.Errorf("SortedKeys[%d] = %q, want %q", i, k, want[i])
		}
	}
}

func TestSort_ByLength(t *testing.T) {
	vars := map[string]string{
		"AB":     "1",
		"ABCDEF": "2",
		"ABC":    "3",
	}
	res, err := Sort(vars, SortOptions{Order: SortByLength})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := []string{"AB", "ABC", "ABCDEF"}
	for i, k := range res.SortedKeys {
		if k != want[i] {
			t.Errorf("SortedKeys[%d] = %q, want %q", i, k, want[i])
		}
	}
}

func TestSort_CaseInsensitive_Ascending(t *testing.T) {
	vars := map[string]string{
		"beta":  "1",
		"ALPHA": "2",
		"gamma": "3",
	}
	res, err := Sort(vars, SortOptions{Order: SortAscending, CaseInsensitive: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := []string{"ALPHA", "beta", "gamma"}
	for i, k := range res.SortedKeys {
		if k != want[i] {
			t.Errorf("SortedKeys[%d] = %q, want %q", i, k, want[i])
		}
	}
}

func TestSort_NilVarsReturnsError(t *testing.T) {
	_, err := Sort(nil, SortOptions{})
	if err == nil {
		t.Fatal("expected error for nil vars, got nil")
	}
}

func TestSort_OriginalUnmodified(t *testing.T) {
	vars := map[string]string{"Z": "1", "A": "2", "M": "3"}
	origCopy := map[string]string{"Z": "1", "A": "2", "M": "3"}
	_, err := Sort(vars, SortOptions{Order: SortDescending})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for k, v := range origCopy {
		if vars[k] != v {
			t.Errorf("original vars[%q] mutated: got %q, want %q", k, vars[k], v)
		}
	}
}
