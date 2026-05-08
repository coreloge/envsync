package envfile

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// Snapshot represents a point-in-time capture of an environment file's key-value pairs.
type Snapshot struct {
	Label     string            `json:"label"`
	Timestamp time.Time         `json:"timestamp"`
	Vars      map[string]string `json:"vars"`
}

// NewSnapshot creates a Snapshot from a parsed env map with a descriptive label.
func NewSnapshot(label string, vars map[string]string) *Snapshot {
	copy := make(map[string]string, len(vars))
	for k, v := range vars {
		copy[k] = v
	}
	return &Snapshot{
		Label:     label,
		Timestamp: time.Now().UTC(),
		Vars:      copy,
	}
}

// SaveSnapshot serializes a Snapshot to a JSON file at the given path.
func SaveSnapshot(path string, s *Snapshot) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("snapshot: create file %q: %w", path, err)
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	if err := enc.Encode(s); err != nil {
		return fmt.Errorf("snapshot: encode: %w", err)
	}
	return nil
}

// LoadSnapshot deserializes a Snapshot from a JSON file at the given path.
func LoadSnapshot(path string) (*Snapshot, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("snapshot: open file %q: %w", path, err)
	}
	defer f.Close()

	var s Snapshot
	if err := json.NewDecoder(f).Decode(&s); err != nil {
		return nil, fmt.Errorf("snapshot: decode: %w", err)
	}
	return &s, nil
}

// DiffSnapshot compares two snapshots and returns a DiffResult using the existing Diff logic.
func DiffSnapshot(base, target *Snapshot) DiffResult {
	return Diff(base.Vars, target.Vars)
}
