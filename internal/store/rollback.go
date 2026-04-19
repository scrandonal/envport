package store

import (
	"fmt"
)

// Rollback restores a snapshot to a previous version from history.
func (m *Manager) Rollback(name string, steps int) error {
	if steps <= 0 {
		steps = 1
	}

	entries, err := m.History(name)
	if err != nil {
		return err
	}
	if len(entries) < steps {
		return fmt.Errorf("not enough history: have %d entries, requested %d steps", len(entries), steps)
	}

	target := entries[len(entries)-steps]

	snap, err := m.Load(target.Name)
	if err != nil {
		return fmt.Errorf("rollback: could not load history entry %q: %w", target.Name, err)
	}

	if err := m.Save(name, snap); err != nil {
		return fmt.Errorf("rollback: could not restore snapshot: %w", err)
	}

	return nil
}
