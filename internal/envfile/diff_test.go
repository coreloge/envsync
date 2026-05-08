package envfile

import (
	"testing"
)

func TestDiff_NoDrift(t *testing.T) {
	source := map[string]string{"FOO": "bar", "BAZ": "qux"}
	target := map[string]string{"FOO": "bar", "BAZ": "qux"}
	result := Diff(source, target)
	if result.HasDrift() {
		t.Errorf("expected no drift, got %+v", result)
	}
}

func TestDiff_MissingInTarget(t *testing.T) {
	source := map[string]string{"FOO": "bar", "MISSING": "val"}
	target := map[string]string{"FOO": "bar"}
	result := Diff(source, target)
	if len(result.MissingInTarget) != 1 || result.MissingInTarget[0] != "MISSING" {
		t.Errorf("expected MISSING in MissingInTarget, got %v", result.MissingInTarget)
	}
}

func TestDiff_ExtraInTarget(t *testing.T) {
	source := map[string]string{"FOO": "bar"}
	target := map[string]string{"FOO": "bar", "EXTRA": "val"}
	result := Diff(source, target)
	if len(result.ExtraInTarget) != 1 || result.ExtraInTarget[0] != "EXTRA" {
		t.Errorf("expected EXTRA in ExtraInTarget, got %v", result.ExtraInTarget)
	}
}

func TestDiff_Mismatched(t *testing.T) {
	source := map[string]string{"FOO": "original"}
	target := map[string]string{"FOO": "changed"}
	result := Diff(source, target)
	pair, ok := result.Mismatched["FOO"]
	if !ok {
		t.Fatal("expected FOO in Mismatched")
	}
	if pair[0] != "original" || pair[1] != "changed" {
		t.Errorf("unexpected mismatch values: %v", pair)
	}
}

func TestDiff_HasDrift(t *testing.T) {
	source := map[string]string{"A": "1"}
	target := map[string]string{"B": "2"}
	result := Diff(source, target)
	if !result.HasDrift() {
		t.Error("expected drift to be detected")
	}
}
