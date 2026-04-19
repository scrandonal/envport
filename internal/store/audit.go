package store

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

type AuditEntry struct {
	Time      time.Time `json:"time"`
	Operation string    `json:"operation"`
	Name      string    `json:"name"`
	Detail    string    `json:"detail,omitempty"`
}

func auditPath(base string) string {
	return filepath.Join(base, "audit.json")
}

func (s *Store) AppendAudit(entry AuditEntry) error {
	entries, _ := s.ReadAudit()
	entries = append(entries, entry)
	data, err := json.MarshalIndent(entries, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(auditPath(s.base), data, 0600)
}

func (s *Store) ReadAudit() ([]AuditEntry, error) {
	data, err := os.ReadFile(auditPath(s.base))
	if err != nil {
		if os.IsNotExist(err) {
			return []AuditEntry{}, nil
		}
		return nil, err
	}
	var entries []AuditEntry
	if err := json.Unmarshal(data, &entries); err != nil {
		return nil, err
	}
	return entries, nil
}

func (s *Store) ClearAudit() error {
	err := os.Remove(auditPath(s.base))
	if os.IsNotExist(err) {
		return nil
	}
	return err
}
