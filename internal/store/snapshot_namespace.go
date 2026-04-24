package store

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

func namespacePath(root, name string) string {
	return filepath.Join(root, name+".namespace.json")
}

// SetNamespace assigns a namespace string to a snapshot.
func SetNamespace(root, name, namespace string) error {
	if _, err := os.Stat(filepath.Join(root, name+".json")); errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("snapshot %q not found", name)
	}
	if namespace == "" {
		return fmt.Errorf("namespace must not be empty")
	}
	data, err := json.Marshal(namespace)
	if err != nil {
		return err
	}
	return os.WriteFile(namespacePath(root, name), data, 0600)
}

// GetNamespace returns the namespace assigned to a snapshot.
// Returns an empty string if none is set.
func GetNamespace(root, name string) (string, error) {
	data, err := os.ReadFile(namespacePath(root, name))
	if errors.Is(err, os.ErrNotExist) {
		return "", nil
	}
	if err != nil {
		return "", err
	}
	var ns string
	if err := json.Unmarshal(data, &ns); err != nil {
		return "", err
	}
	return ns, nil
}

// ClearNamespace removes the namespace assignment from a snapshot.
func ClearNamespace(root, name string) error {
	err := os.Remove(namespacePath(root, name))
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	return err
}

// ListByNamespace returns all snapshot names that belong to the given namespace.
func ListByNamespace(root, namespace string) ([]string, error) {
	entries, err := os.ReadDir(root)
	if err != nil {
		return nil, err
	}
	var results []string
	for _, e := range entries {
		if filepath.Ext(e.Name()) != ".json" {
			continue
		}
		base := e.Name()[:len(e.Name())-len(".json")]
		if filepath.Ext(base) != "" {
			continue // skip metadata files like .namespace.json
		}
		ns, err := GetNamespace(root, base)
		if err != nil || ns != namespace {
			continue
		}
		results = append(results, base)
	}
	return results, nil
}
