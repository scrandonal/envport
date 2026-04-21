package store

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

// OwnerInfo holds ownership metadata for a snapshot.
type OwnerInfo struct {
	User  string `json:"user"`
	Email string `json:"email,omitempty"`
	Team  string `json:"team,omitempty"`
}

func ownerPath(base, name string) string {
	return filepath.Join(base, name+".owner.json")
}

// SetOwner assigns ownership metadata to a named snapshot.
func SetOwner(base, name string, info OwnerInfo) error {
	if !snapshotExists(base, name) {
		return fmt.Errorf("snapshot %q not found", name)
	}
	if info.User == "" {
		return errors.New("owner user must not be empty")
	}
	data, err := json.MarshalIndent(info, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(ownerPath(base, name), data, 0600)
}

// GetOwner returns the ownership metadata for a named snapshot.
// Returns a zero-value OwnerInfo and no error if no owner has been set.
func GetOwner(base, name string) (OwnerInfo, error) {
	data, err := os.ReadFile(ownerPath(base, name))
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return OwnerInfo{}, nil
		}
		return OwnerInfo{}, err
	}
	var info OwnerInfo
	if err := json.Unmarshal(data, &info); err != nil {
		return OwnerInfo{}, err
	}
	return info, nil
}

// ClearOwner removes ownership metadata for a named snapshot.
func ClearOwner(base, name string) error {
	err := os.Remove(ownerPath(base, name))
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}
	return nil
}
