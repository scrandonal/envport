package store

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

var validScopes = map[string]bool{
	"global":  true,
	"local":   true,
	"session": true,
	"user":    true,
}

func scopePath(root, name string) string {
	return filepath.Join(root, name+".scope.json")
}

func SetScope(root, name, scope string) error {
	if !validScopes[scope] {
		return fmt.Errorf("invalid scope %q: must be one of global, local, session, user", scope)
	}
	if _, err := os.Stat(filepath.Join(root, name+".json")); errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("snapshot %q not found", name)
	}
	data, err := json.Marshal(scope)
	if err != nil {
		return err
	}
	return os.WriteFile(scopePath(root, name), data, 0600)
}

func GetScope(root, name string) (string, error) {
	data, err := os.ReadFile(scopePath(root, name))
	if errors.Is(err, os.ErrNotExist) {
		return "", nil
	}
	if err != nil {
		return "", err
	}
	var scope string
	if err := json.Unmarshal(data, &scope); err != nil {
		return "", err
	}
	return scope, nil
}

func ClearScope(root, name string) error {
	err := os.Remove(scopePath(root, name))
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	return err
}

func ListByScope(root, scope string) ([]string, error) {
	entries, err := os.ReadDir(root)
	if err != nil {
		return nil, err
	}
	var names []string
	for _, e := range entries {
		if filepath.Ext(e.Name()) != ".json" {
			continue
		}
		base := e.Name()[:len(e.Name())-5]
		if filepath.Ext(base) != "" {
			continue
		}
		s, err := GetScope(root, base)
		if err != nil {
			continue
		}
		if s == scope {
			names = append(names, base)
		}
	}
	return names, nil
}
