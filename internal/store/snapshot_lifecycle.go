package store

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Lifecycle represents a named stage in a snapshot's lifecycle.
type Lifecycle struct {
	Stage     string    `json:"stage"`
	UpdatedAt time.Time `json:"updated_at"`
}

var validLifecycleStages = map[string]bool{
	"draft":      true,
	"active":     true,
	"deprecated": true,
	"archived":   true,
	"retired":    true,
}

func lifecyclePath(root, name string) string {
	return filepath.Join(root, name+".lifecycle.json")
}

// SetLifecycle sets the lifecycle stage for a snapshot.
func SetLifecycle(root, name, stage string) error {
	if !validLifecycleStages[stage] {
		return fmt.Errorf("invalid lifecycle stage %q: must be one of draft, active, deprecated, archived, retired", stage)
	}
	snapshotFile := filepath.Join(root, name+".json")
	if _, err := os.Stat(snapshotFile); errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("snapshot %q not found", name)
	}
	lc := Lifecycle{Stage: stage, UpdatedAt: time.Now().UTC()}
	data, err := json.MarshalIndent(lc, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(lifecyclePath(root, name), data, 0600)
}

// GetLifecycle returns the lifecycle stage for a snapshot.
func GetLifecycle(root, name string) (Lifecycle, error) {
	data, err := os.ReadFile(lifecyclePath(root, name))
	if errors.Is(err, os.ErrNotExist) {
		return Lifecycle{}, nil
	}
	if err != nil {
		return Lifecycle{}, err
	}
	var lc Lifecycle
	if err := json.Unmarshal(data, &lc); err != nil {
		return Lifecycle{}, err
	}
	return lc, nil
}

// ClearLifecycle removes the lifecycle metadata for a snapshot.
func ClearLifecycle(root, name string) error {
	err := os.Remove(lifecyclePath(root, name))
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	return err
}

// ListByLifecycle returns all snapshot names with the given lifecycle stage.
func ListByLifecycle(root, stage string) ([]string, error) {
	entries, err := os.ReadDir(root)
	if err != nil {
		return nil, err
	}
	var names []string
	for _, e := range entries {
		const suffix = ".lifecycle.json"
		if !e.IsDir() && len(e.Name()) > len(suffix) && e.Name()[len(e.Name())-len(suffix):] == suffix {
			snapshotName := e.Name()[:len(e.Name())-len(suffix)]
			lc, err := GetLifecycle(root, snapshotName)
			if err == nil && lc.Stage == stage {
				names = append(names, snapshotName)
			}
		}
	}
	return names, nil
}
