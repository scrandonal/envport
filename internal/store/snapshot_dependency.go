package store

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

func dependencyPath(base, name string) string {
	return filepath.Join(base, name+".deps.json")
}

// SetDependencies records which snapshots the named snapshot depends on.
func SetDependencies(base, name string, deps []string) error {
	if _, err := os.Stat(filepath.Join(base, name+".json")); errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("snapshot %q not found", name)
	}
	data, err := json.Marshal(deps)
	if err != nil {
		return err
	}
	return os.WriteFile(dependencyPath(base, name), data, 0600)
}

// GetDependencies returns the dependency list for the named snapshot.
func GetDependencies(base, name string) ([]string, error) {
	data, err := os.ReadFile(dependencyPath(base, name))
	if errors.Is(err, os.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var deps []string
	if err := json.Unmarshal(data, &deps); err != nil {
		return nil, err
	}
	return deps, nil
}

// ClearDependencies removes the dependency file for the named snapshot.
func ClearDependencies(base, name string) error {
	err := os.Remove(dependencyPath(base, name))
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	return err
}

// ListDependents returns all snapshots that list name as a dependency.
func ListDependents(base, name string) ([]string, error) {
	entries, err := os.ReadDir(base)
	if err != nil {
		return nil, err
	}
	var dependents []string
	for _, e := range entries {
		if filepath.Ext(e.Name()) != ".json" {
			continue
		}
		candidate := e.Name()[:len(e.Name())-len(".json")]
		if filepath.Ext(candidate) == ".deps" {
			continue
		}
		deps, err := GetDependencies(base, candidate)
		if err != nil {
			continue
		}
		for _, d := range deps {
			if d == name {
				dependents = append(dependents, candidate)
				break
			}
		}
	}
	return dependents, nil
}
