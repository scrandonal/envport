package store

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

var ErrAliasNotFound = errors.New("alias not found")
var ErrAliasExists = errors.New("alias already exists")

func aliasPath(base string) string {
	return filepath.Join(base, "aliases.json")
}

func loadAliases(base string) (map[string]string, error) {
	path := aliasPath(base)
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return map[string]string{}, nil
	}
	if err != nil {
		return nil, err
	}
	var m map[string]string
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, err
	}
	return m, nil
}

func saveAliases(base string, m map[string]string) error {
	data, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(aliasPath(base), data, 0600)
}

func (s *Store) SetAlias(alias, name string, overwrite bool) error {
	if _, err := s.Load(name); err != nil {
		return ErrNotFound
	}
	m, err := loadAliases(s.base)
	if err != nil {
		return err
	}
	if _, exists := m[alias]; exists && !overwrite {
		return ErrAliasExists
	}
	m[alias] = name
	return saveAliases(s.base, m)
}

func (s *Store) ResolveAlias(alias string) (string, error) {
	m, err := loadAliases(s.base)
	if err != nil {
		return "", err
	}
	name, ok := m[alias]
	if !ok {
		return "", ErrAliasNotFound
	}
	return name, nil
}

func (s *Store) DeleteAlias(alias string) error {
	m, err := loadAliases(s.base)
	if err != nil {
		return err
	}
	if _, ok := m[alias]; !ok {
		return ErrAliasNotFound
	}
	delete(m, alias)
	return saveAliases(s.base, m)
}

func (s *Store) ListAliases() (map[string]string, error) {
	return loadAliases(s.base)
}
