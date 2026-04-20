package store

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

// Priority levels for snapshots.
const (
	PriorityLow    = "low"
	PriorityNormal = "normal"
	PriorityHigh   = "high"
	PriorityCritical = "critical"
)

var validPriorities = map[string]bool{
	PriorityLow:      true,
	PriorityNormal:   true,
	PriorityHigh:     true,
	PriorityCritical: true,
}

type PriorityRecord struct {
	Level string `json:"level"`
}

func priorityPath(root, name string) string {
	return filepath.Join(root, name+".priority.json")
}

// SetPriority assigns a priority level to a named snapshot.
func SetPriority(root, name, level string) error {
	if !validPriorities[level] {
		return fmt.Errorf("invalid priority %q: must be one of low, normal, high, critical", level)
	}
	snapPath := filepath.Join(root, name+".json")
	if _, err := os.Stat(snapPath); errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("snapshot %q not found", name)
	}
	rec := PriorityRecord{Level: level}
	data, err := json.Marshal(rec)
	if err != nil {
		return err
	}
	return os.WriteFile(priorityPath(root, name), data, 0600)
}

// GetPriority returns the priority level for a named snapshot.
// Returns PriorityNormal if none is set.
func GetPriority(root, name string) (string, error) {
	data, err := os.ReadFile(priorityPath(root, name))
	if errors.Is(err, os.ErrNotExist) {
		return PriorityNormal, nil
	}
	if err != nil {
		return "", err
	}
	var rec PriorityRecord
	if err := json.Unmarshal(data, &rec); err != nil {
		return "", err
	}
	return rec.Level, nil
}

// ClearPriority removes the priority record for a named snapshot.
func ClearPriority(root, name string) error {
	err := os.Remove(priorityPath(root, name))
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	return err
}

// ListByPriority returns all snapshot names that have the given priority level.
func ListByPriority(root, level string) ([]string, error) {
	entries, err := os.ReadDir(root)
	if err != nil {
		return nil, err
	}
	var names []string
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		const suffix = ".priority.json"
		if len(e.Name()) <= len(suffix) || e.Name()[len(e.Name())-len(suffix):] != suffix {
			continue
		}
		name := e.Name()[:len(e.Name())-len(suffix)]
		l, err := GetPriority(root, name)
		if err != nil {
			continue
		}
		if l == level {
			names = append(names, name)
		}
	}
	return names, nil
}
