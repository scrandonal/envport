package store

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

// EnvironmentRecord stores the environment context (e.g. hostname, user, OS)
// captured at the time a snapshot was saved.
type EnvironmentRecord struct {
	Hostname string `json:"hostname"`
	User     string `json:"user"`
	OS       string `json:"os"`
	Shell    string `json:"shell"`
}

func environmentPath(base, name string) string {
	return filepath.Join(base, name+".env_context.json")
}

// SetEnvironment persists the environment record for the named snapshot.
func SetEnvironment(base, name string, rec EnvironmentRecord) error {
	if !snapshotExists(base, name) {
		return fmt.Errorf("snapshot %q not found", name)
	}
	data, err := json.MarshalIndent(rec, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(environmentPath(base, name), data, 0600)
}

// GetEnvironment loads the environment record for the named snapshot.
func GetEnvironment(base, name string) (EnvironmentRecord, error) {
	data, err := os.ReadFile(environmentPath(base, name))
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return EnvironmentRecord{}, nil
		}
		return EnvironmentRecord{}, err
	}
	var rec EnvironmentRecord
	if err := json.Unmarshal(data, &rec); err != nil {
		return EnvironmentRecord{}, err
	}
	return rec, nil
}

// ClearEnvironment removes the environment record for the named snapshot.
func ClearEnvironment(base, name string) error {
	err := os.Remove(environmentPath(base, name))
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	return err
}
