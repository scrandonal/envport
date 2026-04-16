package store

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

const defaultDirName = ".envport"

// Store manages named snapshot files on disk.
type Store struct {
	Dir string
}

// Default returns a Store rooted in the user's home directory.
func Default() (*Store, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	return &Store{Dir: filepath.Join(home, defaultDirName)}, nil
}

// New creates a Store at the given directory.
func New(dir string) *Store {
	return &Store{Dir: dir}
}

// Init ensures the store directory exists with safe permissions.
func (s *Store) Init() error {
	return os.MkdirAll(s.Dir, 0700)
}

// Path returns the file path for a named snapshot.
func (s *Store) Path(name string) string {
	return filepath.Join(s.Dir, name+".json")
}

// List returns all snapshot names in the store.
func (s *Store) List() ([]string, error) {
	entries, err := os.ReadDir(s.Dir)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, nil
		}
		return nil, err
	}
	var names []string
	for _, e := range entries {
		if !e.IsDir() && filepath.Ext(e.Name()) == ".json" {
			names = append(names, e.Name()[:len(e.Name())-5])
		}
	}
	return names, nil
}

// Delete removes a named snapshot from the store.
func (s *Store) Delete(name string) error {
	err := os.Remove(s.Path(name))
	if errors.Is(err, os.ErrNotExist) {
		return ErrNotFound
	}
	return err
}

// Exists reports whether a named snapshot exists.
func (s *Store) Exists(name string) bool {
	_, err := os.Stat(s.Path(name))
	return err == nil
}

// ErrNotFound is returned when a snapshot name does not exist in the store.
var ErrNotFound = errors.New("snapshot not found")

// metadata is used internally to peek at a snapshot file's name field.
type metadata struct {
	Name string `json:"name"`
}

// ReadName reads only the name field from a snapshot file.
func (s *Store) ReadName(name string) (string, error) {
	data, err := os.ReadFile(s.Path(name))
	if err != nil {
		return "", err
	}
	var m metadata
	if err := json.Unmarshal(data, &m); err != nil {
		return "", err
	}
	return m.Name, nil
}
