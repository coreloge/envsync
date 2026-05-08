package envfile

import "fmt"

// AuditSeverity represents the severity level of an audit finding.
type AuditSeverity string

const (
	SeverityInfo    AuditSeverity = "info"
	SeverityWarning AuditSeverity = "warning"
	SeverityCritical AuditSeverity = "critical"
)

// AuditFinding represents a single finding from an environment audit.
type AuditFinding struct {
	Key      string
	Message  string
	Severity AuditSeverity
}

func (f AuditFinding) String() string {
	return fmt.Sprintf("[%s] %s: %s", f.Severity, f.Key, f.Message)
}

// AuditResult holds the complete result of an audit operation.
type AuditResult struct {
	Findings []AuditFinding
}

// HasIssues returns true if there are any findings.
func (r *AuditResult) HasIssues() bool {
	return len(r.Findings) > 0
}

// BySeverity returns findings filtered by the given severity.
func (r *AuditResult) BySeverity(s AuditSeverity) []AuditFinding {
	var out []AuditFinding
	for _, f := range r.Findings {
		if f.Severity == s {
			out = append(out, f)
		}
	}
	return out
}

// Audit compares a source env map against a target env map and produces
// an AuditResult describing drift, missing keys, and empty values.
func Audit(source, target map[string]string) *AuditResult {
	result := &AuditResult{}

	for key, srcVal := range source {
		tgtVal, exists := target[key]
		if !exists {
			result.Findings = append(result.Findings, AuditFinding{
				Key:      key,
				Message:  "key present in source but missing in target",
				Severity: SeverityCritical,
			})
			continue
		}
		if srcVal != tgtVal {
			result.Findings = append(result.Findings, AuditFinding{
				Key:      key,
				Message:  "value differs between source and target",
				Severity: SeverityWarning,
			})
		}
		if tgtVal == "" {
			result.Findings = append(result.Findings, AuditFinding{
				Key:      key,
				Message:  "value is empty in target",
				Severity: SeverityWarning,
			})
		}
	}

	for key := range target {
		if _, exists := source[key]; !exists {
			result.Findings = append(result.Findings, AuditFinding{
				Key:      key,
				Message:  "key exists in target but not in source",
				Severity: SeverityInfo,
			})
		}
	}

	return result
}
