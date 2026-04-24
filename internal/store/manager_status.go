package store

import "fmt"

func (m *Manager) SetStatus(name, status string) error {
	if err := m.requireExists(name); err != nil {
		return err
	}
	return SetStatus(m.store.Base(), name, status)
}

func (m *Manager) GetStatus(name string) (string, error) {
	return GetStatus(m.store.Base(), name)
}

func (m *Manager) ClearStatus(name string) error {
	return ClearStatus(m.store.Base(), name)
}

func (m *Manager) ListByStatus(status string) ([]string, error) {
	if !validStatuses[status] {
		return nil, fmt.Errorf("invalid status %q", status)
	}
	return ListByStatus(m.store.Base(), status)
}
