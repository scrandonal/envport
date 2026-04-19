package store

import "time"

func (m *Manager) RecordAudit(operation, name, detail string) error {
	return m.store.AppendAudit(AuditEntry{
		Time:      time.Now(),
		Operation: operation,
		Name:      name,
		Detail:    detail,
	})
}

func (m *Manager) AuditLog() ([]AuditEntry, error) {
	return m.store.ReadAudit()
}

func (m *Manager) ClearAuditLog() error {
	return m.store.ClearAudit()
}
