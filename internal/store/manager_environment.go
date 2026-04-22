package store

// SetEnvironment records the capture-time environment context for a snapshot.
func (m *Manager) SetEnvironment(name string, rec EnvironmentRecord) error {
	return SetEnvironment(m.store.Base(), name, rec)
}

// GetEnvironment returns the environment context stored for a snapshot.
func (m *Manager) GetEnvironment(name string) (EnvironmentRecord, error) {
	return GetEnvironment(m.store.Base(), name)
}

// ClearEnvironment removes the environment context for a snapshot.
func (m *Manager) ClearEnvironment(name string) error {
	return ClearEnvironment(m.store.Base(), name)
}
