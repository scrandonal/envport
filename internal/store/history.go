package store

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

// HistoryEntry records a snapshot operation for auditing.
type HistoryEntry struct {
	Name      string    `json:"name"`
	Operation string    `json:"operation"`
	Timestamp time.Time `json:"timestamp"`
}

const historyFile = "history.json"

// AppendHistory adds an entry to the store's history log.
func (s *Store) AppendHistory(name, operation string) error {
	entries, _ := s.ReadHistory()
	entries = append(entries, HistoryEntry{
		Name:      name,
		Operation: operation,
		Timestamp: time.Now().UTC(),
	})
	data, err := json.MarshalIndent(entries, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(s.dir, historyFile), data, 0600)
}

// ReadHistory returns all recorded history entries.
func (s *Store) ReadHistory() ([]HistoryEntry, error) {
	path := filepath.Join(s.dir, historyFile)
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return []HistoryEntry{}, nil
		}
		return nil, err
	}
	var entries []HistoryEntry
	if err := json.Unmarshal(data, &entries); err != nil {
		return nil, err
	}
	return entries, nil
}

// ClearHistory removes all history entries.
func (s *Store) ClearHistory() error {
	path := filepath.Join(s.dir, historyFile)
	err := os.Remove(path)
	if os.IsNotExist(err) {
		return nil
	}
	return err
}
