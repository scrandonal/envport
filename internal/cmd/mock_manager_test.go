package cmd

import (
	"errors"
	"fmt"

	"github.com/user/envport/internal/snapshot"
)

type mockManager struct {
	snapshots map[string]map[string]string
}

func (m *mockManager) init() {
	if m.snapshots == nil {
		m.snapshots = make(map[string]map[string]string)
	}
}

func (m *mockManager) Save(name string, snap *snapshot.Snapshot) error {
	m.init()
	m.snapshots[name] = snap.Vars
	return nil
}

func (m *mockManager) Load(name string) (*snapshot.Snapshot, error) {
	m.init()
	vars, ok := m.snapshots[name]
	if !ok {
		return nil, fmt.Errorf("snapshot %q not found", name)
	}
	return &snapshot.Snapshot{Vars: vars}, nil
}

func (m *mockManager) List() ([]string, error) {
	m.init()
	names := make([]string, 0, len(m.snapshots))
	for k := range m.snapshots {
		names = append(names, k)
	}
	return names, nil
}

func (m *mockManager) Delete(name string) error {
	m.init()
	if _, ok := m.snapshots[name]; !ok {
		return errors.New("not found")
	}
	delete(m.snapshots, name)
	return nil
}

func (m *mockManager) Rename(oldName, newName string) error {
	m.init()
	vars, ok := m.snapshots[oldName]
	if !ok {
		return fmt.Errorf("snapshot %q not found", oldName)
	}
	m.snapshots[newName] = vars
	delete(m.snapshots, oldName)
	return nil
}
