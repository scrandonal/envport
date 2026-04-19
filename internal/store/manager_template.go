package store

import "fmt"

func (m *Manager) SaveTemplate(t Template) error {
	return m.store.SaveTemplate(t)
}

func (m *Manager) LoadTemplate(name string) (Template, error) {
	return m.store.LoadTemplate(name)
}

func (m *Manager) DeleteTemplate(name string) error {
	return m.store.DeleteTemplate(name)
}

func (m *Manager) ListTemplates() ([]string, error) {
	return m.store.ListTemplates()
}

// ApplyTemplate creates a new snapshot from a template, filling in defaults
// and leaving other keys empty. Returns error if snapshot already exists.
func (m *Manager) ApplyTemplate(templateName, snapshotName string) error {
	tmpl, err := m.store.LoadTemplate(templateName)
	if err != nil {
		return fmt.Errorf("load template %q: %w", templateName, err)
	}
	if _, err := m.store.Load(snapshotName); err == nil {
		return fmt.Errorf("snapshot %q already exists", snapshotName)
	}
	vars := make(map[string]string, len(tmpl.Keys))
	for _, k := range tmpl.Keys {
		if v, ok := tmpl.Defaults[k]; ok {
			vars[k] = v
		} else {
			vars[k] = ""
		}
	}
	return m.store.Save(snapshotName, vars)
}
