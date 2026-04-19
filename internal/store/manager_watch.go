package store

func (m *Manager) SetWatch(name string) error {
	return SetWatch(m.store.Root, name)
}

func (m *Manager) GetWatch(name string) (*WatchEvent, error) {
	return GetWatch(m.store.Root, name)
}

func (m *Manager) ClearWatch(name string) error {
	return ClearWatch(m.store.Root, name)
}

func (m *Manager) ListWatched() ([]WatchEvent, error) {
	return ListWatched(m.store.Root)
}
