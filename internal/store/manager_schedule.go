package store

func (m *Manager) SetSchedule(name string, s Schedule) error {
	return SetSchedule(m.store.Base(), name, s)
}

func (m *Manager) GetSchedule(name string) (Schedule, error) {
	return GetSchedule(m.store.Base(), name)
}

func (m *Manager) ClearSchedule(name string) error {
	return ClearSchedule(m.store.Base(), name)
}

func (m *Manager) ListScheduled() ([]string, error) {
	return ListScheduled(m.store.Base())
}
