package envfile

import (
	"testing"
)

func TestAudit_NoFindings(t *testing.T) {
	src := map[string]string{"FOO": "bar", "BAZ": "qux"}
	tgt := map[string]string{"FOO": "bar", "BAZ": "qux"}

	result := Audit(src, tgt)
	if result.HasIssues() {
		t.Errorf("expected no findings, got %d", len(result.Findings))
	}
}

func TestAudit_MissingInTarget(t *testing.T) {
	src := map[string]string{"FOO": "bar", "MISSING": "val"}
	tgt := map[string]string{"FOO": "bar"}

	result := Audit(src, tgt)
	critical := result.BySeverity(SeverityCritical)
	if len(critical) != 1 {
		t.Fatalf("expected 1 critical finding, got %d", len(critical))
	}
	if critical[0].Key != "MISSING" {
		t.Errorf("expected key MISSING, got %s", critical[0].Key)
	}
}

func TestAudit_ExtraInTarget(t *testing.T) {
	src := map[string]string{"FOO": "bar"}
	tgt := map[string]string{"FOO": "bar", "EXTRA": "val"}

	result := Audit(src, tgt)
	info := result.BySeverity(SeverityInfo)
	if len(info) != 1 {
		t.Fatalf("expected 1 info finding, got %d", len(info))
	}
	if info[0].Key != "EXTRA" {
		t.Errorf("expected key EXTRA, got %s", info[0].Key)
	}
}

func TestAudit_ValueMismatch(t *testing.T) {
	src := map[string]string{"FOO": "original"}
	tgt := map[string]string{"FOO": "changed"}

	result := Audit(src, tgt)
	warnings := result.BySeverity(SeverityWarning)
	if len(warnings) != 1 {
		t.Fatalf("expected 1 warning, got %d", len(warnings))
	}
	if warnings[0].Key != "FOO" {
		t.Errorf("expected key FOO, got %s", warnings[0].Key)
	}
}

func TestAudit_EmptyValueInTarget(t *testing.T) {
	src := map[string]string{"DB_PASS": "secret"}
	tgt := map[string]string{"DB_PASS": ""}

	result := Audit(src, tgt)
	// mismatch + empty both produce warnings
	warnings := result.BySeverity(SeverityWarning)
	if len(warnings) < 1 {
		t.Fatalf("expected at least 1 warning for empty value, got %d", len(warnings))
	}
}

func TestAuditFinding_String(t *testing.T) {
	f := AuditFinding{Key: "FOO", Message: "some issue", Severity: SeverityCritical}
	got := f.String()
	expected := "[critical] FOO: some issue"
	if got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}
}
