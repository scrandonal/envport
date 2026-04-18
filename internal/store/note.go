package store

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

var ErrNoteNotFound = errors.New("note not found")

func notePath(base, name string) string {
	return filepath.Join(base, name+".note.json")
}

func (s *Store) SetNote(name, text string) error {
	if !s.Exists(name) {
		return ErrNotFound
	}
	p := notePath(s.dir, name)
	data, err := json.Marshal(text)
	if err != nil {
		return err
	}
	return os.WriteFile(p, data, 0600)
}

func (s *Store) GetNote(name string) (string, error) {
	if !s.Exists(name) {
		return "", ErrNotFound
	}
	p := notePath(s.dir, name)
	data, err := os.ReadFile(p)
	if errors.Is(err, os.ErrNotExist) {
		return "", ErrNoteNotFound
	}
	if err != nil {
		return "", err
	}
	var text string
	if err := json.Unmarshal(data, &text); err != nil {
		return "", err
	}
	return text, nil
}

func (s *Store) ClearNote(name string) error {
	p := notePath(s.dir, name)
	err := os.Remove(p)
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	return err
}
