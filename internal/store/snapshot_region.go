package store

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

func regionPath(base, name string) string {
	return filepath.Join(base, name+".region.json")
}

// SetRegion associates a region string (e.g. "us-east-1") with a snapshot.
func SetRegion(base, name, region string) error {
	if !snapshotExists(base, name) {
		return fmt.Errorf("snapshot %q not found", name)
	}
	if region == "" {
		return errors.New("region must not be empty")
	}
	data, err := json.Marshal(region)
	if err != nil {
		return err
	}
	return os.WriteFile(regionPath(base, name), data, 0600)
}

// GetRegion returns the region associated with a snapshot.
// Returns ("", nil) if no region has been set.
func GetRegion(base, name string) (string, error) {
	data, err := os.ReadFile(regionPath(base, name))
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return "", nil
		}
		return "", err
	}
	var region string
	if err := json.Unmarshal(data, &region); err != nil {
		return "", err
	}
	return region, nil
}

// ClearRegion removes the region metadata for a snapshot.
func ClearRegion(base, name string) error {
	err := os.Remove(regionPath(base, name))
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}
	return nil
}

// ListByRegion returns snapshot names that match the given region.
func ListByRegion(base, region string) ([]string, error) {
	entries, err := os.ReadDir(base)
	if err != nil {
		return nil, err
	}
	var matches []string
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := snapshotNameFromEntry(e.Name())
		if name == "" {
			continue
		}
		r, err := GetRegion(base, name)
		if err != nil {
			continue
		}
		if r == region {
			matches = append(matches, name)
		}
	}
	return matches, nil
}
