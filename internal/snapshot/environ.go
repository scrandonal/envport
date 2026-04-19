package snapshot

import (
	"os"
	"strings"
)

// OSEnviron returns the current process environment as a map.
func OSEnviron() map[string]string {
	return ParseEnviron(os.Environ())
}

// ParseEnviron converts a slice of "KEY=VALUE" strings to a map.
func ParseEnviron(pairs []string) map[string]string {
	m := make(map[string]string, len(pairs))
	for _, p := range pairs {
		parts := strings.SplitN(p, "=", 2)
		if len(parts) == 2 {
			m[parts[0]] = parts[1]
		}
	}
	return m
}
