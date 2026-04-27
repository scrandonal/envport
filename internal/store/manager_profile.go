package store

func (m *Manager) SetProfile(name string, p Profile) error {
	return SetProfile(m.store.Base(), name, p)
}

func (m *Manager) GetProfile(name string) (Profile, error) {
	return GetProfile(m.store.Base(), name)
}

func (m *Manager) ClearProfile(name string) error {
	return ClearProfile(m.store.Base(), name)
}
