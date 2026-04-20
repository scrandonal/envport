package store

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

// Visibility represents the visibility level of a snapshot.
type Visibility string

const (
	VisibilityPrivate Visibility = "private"
	VisibilityShared  Visibility = "shared"
	VisibilityPublic  Visibility = "public"
)

var validVisibilities = map[Visibility]bool{
	VisibilityPrivate: true,
	VisibilityShared:  true,
	VisibilityPublic:  true,
}

func visibilityPath(base, name string) string {
	return filepath.Join(base, name+".visibility.json")
}

// SetVisibility sets the visibility level for a named snapshot.
func SetVisibility(base, name string, v Visibility) error {
	if !validVisibilities[v] {
		return fmt.Errorf("invalid visibility %q: must be private, shared, or public", v)
	}
	if _, err := os.Stat(filepath.Join(base, name+".json")); errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("snapshot %q not found", name)
	}
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}
	return os.WriteFile(visibilityPath(base, name), data, 0600)
}

// GetVisibility returns the visibility level for a named snapshot.
// Defaults to VisibilityPrivate if not set.
func GetVisibility(base, name string) (Visibility, error) {
	data, err := os.ReadFile(visibilityPath(base, name))
	if errors.Is(err, os.ErrNotExist) {
		return VisibilityPrivate, nil
	}
	if err != nil {
		return "", err
	}
	var v Visibility
	if err := json.Unmarshal(data, &v); err != nil {
		return "", err
	}
	return v, nil
}

// ClearVisibility removes the visibility setting for a named snapshot.
func ClearVisibility(base, name string) error {
	err := os.Remove(visibilityPath(base, name))
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	return err
}

// ListByVisibility returns snapshot names that match the given visibility level.
func ListByVisibility(base string, v Visibility, all []string) ([]string, error) {
	var result []string
	for _, name := range all {
		current, err := GetVisibility(base, name)
		if err != nil {
			return nil, err
		}
		if current == v {
			result = append(result, name)
		}
	}
	return result, nil
}
