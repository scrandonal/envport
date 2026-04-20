package store

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

func categoryPath(base, name string) string {
	return filepath.Join(base, name+".category.json")
}

// SetCategory assigns a category string to a named snapshot.
func SetCategory(base, name, category string) error {
	if _, err := os.Stat(filepath.Join(base, name+".json")); errors.Is(err, os.ErrNotExist) {
		return ErrNotFound
	}
	data, err := json.Marshal(category)
	if err != nil {
		return err
	}
	return os.WriteFile(categoryPath(base, name), data, 0600)
}

// GetCategory returns the category assigned to a named snapshot.
// Returns an empty string and no error if no category is set.
func GetCategory(base, name string) (string, error) {
	data, err := os.ReadFile(categoryPath(base, name))
	if errors.Is(err, os.ErrNotExist) {
		return "", nil
	}
	if err != nil {
		return "", err
	}
	var category string
	if err := json.Unmarshal(data, &category); err != nil {
		return "", err
	}
	return category, nil
}

// ClearCategory removes the category assignment for a named snapshot.
func ClearCategory(base, name string) error {
	err := os.Remove(categoryPath(base, name))
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	return err
}

// ListByCategory returns all snapshot names that belong to the given category.
func ListByCategory(base, category string) ([]string, error) {
	entries, err := os.ReadDir(base)
	if errors.Is(err, os.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var names []string
	for _, e := range entries {
		if filepath.Ext(e.Name()) != ".json" {
			continue
		}
		base2 := e.Name()[:len(e.Name())-len(".json")]
		if filepath.Ext(base2) != "" {
			continue // skip meta files like .category.json
		}
		cat, err := GetCategory(base, base2)
		if err != nil {
			continue
		}
		if cat == category {
			names = append(names, base2)
		}
	}
	return names, nil
}
