package store

import (
	"fmt"

	"github.com/user/envport/internal/snapshot"
)

// Manager combines a Store with snapshot operations.
type Manager struct {
	store *Store
}

// NewManager creates a Manager backed by the given Store.
func NewManager(s *Store) *Manager {
	return &Manager{store: s}
}

// Save captures env vars into a named snapshot and persists it.
func (m *Manager) Save(name string, env []string) error {
	if err := m.store.Init(); err != nil {
		return fmt.Errorf("init store: %w", err)
	}
	snap := snapshot.FromEnviron(name, env)
	if err := snap.Save(m.store.Path(name)); err != nil {
		return fmt.Errorf("save snapshot %q: %w", name, err)
	}
	return nil
}

// Load reads a named snapshot from the store.
func (m *Manager) Load(name string) (*snapshot.Snapshot, error) {
	if !m.store.Exists(name) {
		return nil, fmt.Errorf("%w: %s", ErrNotFound, name)
	}
	snap, err := snapshot.Load(m.store.Path(name))
	if err != nil {
		return nil, fmt.Errorf("load snapshot %q: %w", name, err)
	}
	return snap, nil
}

// List returns the names of all stored snapshots.
func (m *Manager) List() ([]string, error) {
	return m.store.List()
}

// Delete removes a named snapshot from the store.
func (m *Manager) Delete(name string) error {
	return m.store.Delete(name)
}

// Rename moves a snapshot from oldName to newName.
func (m *Manager) Rename(oldName, newName string) error {
	if !m.store.Exists(oldName) {
		return fmt.Errorf("%w: %s", ErrNotFound, oldName)
	}
	if m.store.Exists(newName) {
		return fmt.Errorf("snapshot %q already exists", newName)
	}
	snap, err := snapshot.Load(m.store.Path(oldName))
	if err != nil {
		return fmt.Errorf("load: %w", err)
	}
	snap.Name = newName
	if err := snap.Save(m.store.Path(newName)); err != nil {
		return fmt.Errorf("save: %w", err)
	}
	return m.store.Delete(oldName)
}
