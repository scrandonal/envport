package store

import "fmt"

// Clone duplicates a snapshot under a new name, optionally with a tag.
func (m *Manager) Clone(src, dst string, overwrite bool) error {
	snap, err := m.Load(src)
	if err != nil {
		return fmt.Errorf("clone: source %q not found", src)
	}

	exists, _ := m.store.Exists(dst)
	if exists && !overwrite {
		return fmt.Errorf("clone: destination %q already exists (use --overwrite to replace)", dst)
	}

	cloned := snap.Clone()
	cloned.Meta.Tags = append([]string{}, snap.Meta.Tags...)

	if err := m.Save(dst, cloned); err != nil {
		return fmt.Errorf("clone: failed to save destination: %w", err)
	}
	return nil
}
