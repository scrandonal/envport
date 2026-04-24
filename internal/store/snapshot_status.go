package store

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

var validStatuses = map[string]bool{
	"active":     true,
	"deprecated": true,
	"draft":      true,
	"archived":   true,
}

type StatusRecord struct {
	Status string `json:"status"`
}

func statusPath(base, name string) string {
	return filepath.Join(base, name+".status.json")
}

func SetStatus(base, name, status string) error {
	if !validStatuses[status] {
		return fmt.Errorf("invalid status %q: must be one of active, deprecated, draft, archived", status)
	}
	if _, err := os.Stat(filepath.Join(base, name+".json")); errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("snapshot %q not found", name)
	}
	rec := StatusRecord{Status: status}
	data, err := json.MarshalIndent(rec, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(statusPath(base, name), data, 0600)
}

func GetStatus(base, name string) (string, error) {
	data, err := os.ReadFile(statusPath(base, name))
	if errors.Is(err, os.ErrNotExist) {
		return "", nil
	}
	if err != nil {
		return "", err
	}
	var rec StatusRecord
	if err := json.Unmarshal(data, &rec); err != nil {
		return "", err
	}
	return rec.Status, nil
}

func ClearStatus(base, name string) error {
	err := os.Remove(statusPath(base, name))
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	return err
}

func ListByStatus(base, status string) ([]string, error) {
	entries, err := os.ReadDir(base)
	if err != nil {
		return nil, err
	}
	var names []string
	for _, e := range entries {
		if filepath.Ext(e.Name()) != ".json" {
			continue
		}
		base2 := e.Name()[:len(e.Name())-5]
		if filepath.Ext(base2) != "" {
			continue
		}
		s, err := GetStatus(base, base2)
		if err != nil {
			continue
		}
		if s == status {
			names = append(names, base2)
		}
	}
	return names, nil
}
