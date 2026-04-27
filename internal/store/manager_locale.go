package store

func (m *Manager) SetLocale(name, locale string) error {
	return SetLocale(m.store.Base(), name, locale)
}

func (m *Manager) GetLocale(name string) (string, error) {
	return GetLocale(m.store.Base(), name)
}

func (m *Manager) ClearLocale(name string) error {
	return ClearLocale(m.store.Base(), name)
}

func (m *Manager) ListByLocale(locale string) ([]string, error) {
	return ListByLocale(m.store.Base(), locale)
}
