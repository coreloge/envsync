package envfile

import (
	"fmt"
	"os"
	"sort"
	"strings"
)

// Write serialises vars to an env file at path.
// Keys are written in sorted order so output is deterministic.
func Write(vars map[string]string, path string) error {
	keys := make([]string, 0, len(vars))
	for k := range vars {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var sb strings.Builder
	for _, k := range keys {
		v := vars[k]
		if needsQuoting(v) {
			fmt.Fprintf(&sb, "%s=\"%s\"\n", k, v)
		} else {
			fmt.Fprintf(&sb, "%s=%s\n", k, v)
		}
	}

	if err := os.WriteFile(path, []byte(sb.String()), 0o644); err != nil {
		return fmt.Errorf("write: %w", err)
	}
	return nil
}

// needsQuoting reports whether a value should be wrapped in double quotes.
// Values containing spaces, tabs, or hash characters require quoting so that
// shells and env-file parsers do not misinterpret them.
func needsQuoting(v string) bool {
	return strings.ContainsAny(v, " \t#=\n")
}
