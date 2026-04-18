package store

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

var ErrNotPinned = errors.New("no snapshot is pinned")

type pinFile struct {
	Name string `json:"name"`
}

func pinPath(dir string) string {
	return filepath.Join(dir, "pin.json")
}

// Pin marks the given snapshot name as the active/pinned snapshot.
func (s *Store) Pin(name string) error {
	if _, err := s.Load(name); err != nil {
		return err
	}
	data, err := json.Marshal(pinFile{Name: name})
	if err != nil {
		return err
	}
	return os.WriteFile(pinPath(s.dir), data, 0600)
}

// Unpin removes the pinned snapshot marker.
func (s *Store) Unpin() error {
	err := os.Remove(pinPath(s.dir))
	if errors.Is(err, os.ErrNotExist) {
		return ErrNotPinned
	}
	return err
}

// Pinned returns the name of the currently pinned snapshot.
func (s *Store) Pinned() (string, error) {
	data, err := os.ReadFile(pinPath(s.dir))
	if errors.Is(err, os.ErrNotExist) {
		return "", ErrNotPinned
	}
	if err != nil {
		return "", err
	}
	var pf pinFile
	if err := json.Unmarshal(data, &pf); err != nil {
		return "", err
	}
	return pf.Name, nil
}
