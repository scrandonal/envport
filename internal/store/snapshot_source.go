package store

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

// SourceRecord describes where a snapshot was originally captured from.
type SourceRecord struct {
	Hostname  string `json:"hostname"`
	Directory string `json:"directory"`
	User      string `json:"user"`
}

func sourcePath(base, name string) string {
	return filepath.Join(base, name+".source.json")
}

// SetSource records the capture source for a named snapshot.
func SetSource(base, name string, rec SourceRecord) error {
	if _, err := os.Stat(filepath.Join(base, name+".json")); errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("snapshot %q not found", name)
	}
	data, err := json.MarshalIndent(rec, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(sourcePath(base, name), data, 0600)
}

// GetSource returns the source record for a named snapshot.
// Returns a zero-value SourceRecord and no error when no record exists.
func GetSource(base, name string) (SourceRecord, error) {
	data, err := os.ReadFile(sourcePath(base, name))
	if errors.Is(err, os.ErrNotExist) {
		return SourceRecord{}, nil
	}
	if err != nil {
		return SourceRecord{}, err
	}
	var rec SourceRecord
	if err := json.Unmarshal(data, &rec); err != nil {
		return SourceRecord{}, err
	}
	return rec, nil
}

// ClearSource removes the source record for a named snapshot.
func ClearSource(base, name string) error {
	err := os.Remove(sourcePath(base, name))
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	return err
}
