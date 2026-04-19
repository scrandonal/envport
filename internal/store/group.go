package store

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

var ErrGroupNotFound = errors.New("group not found")
var ErrGroupExists = errors.New("group already exists")

func groupPath(base, name string) string {
	return filepath.Join(base, "groups", name+".json")
}

func (s *Store) CreateGroup(name string, snapshots []string) error {
	for _, sn := range snapshots {
		if !s.Exists(sn) {
			return ErrNotFound
		}
	}
	p := groupPath(s.base, name)
	if _, err := os.Stat(p); err == nil {
		return ErrGroupExists
	}
	if err := os.MkdirAll(filepath.Dir(p), 0700); err != nil {
		return err
	}
	data, err := json.Marshal(snapshots)
	if err != nil {
		return err
	}
	return os.WriteFile(p, data, 0600)
}

func (s *Store) GetGroup(name string) ([]string, error) {
	p := groupPath(s.base, name)
	data, err := os.ReadFile(p)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, ErrGroupNotFound
		}
		return nil, err
	}
	var snapshots []string
	if err := json.Unmarshal(data, &snapshots); err != nil {
		return nil, err
	}
	return snapshots, nil
}

func (s *Store) DeleteGroup(name string) error {
	p := groupPath(s.base, name)
	if _, err := os.Stat(p); errors.Is(err, os.ErrNotExist) {
		return ErrGroupNotFound
	}
	return os.Remove(p)
}

func (s *Store) ListGroups() ([]string, error) {
	dir := filepath.Join(s.base, "groups")
	entries, err := os.ReadDir(dir)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return []string{}, nil
		}
		return nil, err
	}
	var names []string
	for _, e := range entries {
		if !e.IsDir() {
			names = append(names, e.Name()[:len(e.Name())-5])
		}
	}
	return names, nil
}
