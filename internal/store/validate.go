package store

import (
	"fmt"
	"regexp"
	"sort"
)

var validKeyRe = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*$`)

// ValidationError holds all invalid keys found in a snapshot.
type ValidationError struct {
	InvalidKeys []string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("snapshot contains invalid keys: %v", e.InvalidKeys)
}

// Validate checks that all keys in the named snapshot are valid shell
// variable names. Returns a *ValidationError listing any bad keys, or
// nil if the snapshot is clean.
func (m *Manager) Validate(name string) error {
	snap, err := m.Load(name)
	if err != nil {
		return err
	}

	var bad []string
	for k := range snap.Vars {
		if !validKeyRe.MatchString(k) {
			bad = append(bad, k)
		}
	}
	if len(bad) > 0 {
		sort.Strings(bad)
		return &ValidationError{InvalidKeys: bad}
	}
	return nil
}
