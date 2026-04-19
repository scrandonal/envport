package store

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"time"
)

// Meta holds optional metadata associated with a snapshot.
type Meta struct {
	Description string    `json:"description,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func metaPath(base, name string) string {
	return filepath.Join(base, name+".meta.json")
}

func (s *Store) SetMeta(name, description string) error {
	if !s.Exists(name) {
		return errors.New("snapshot not found: " + name)
	}
	p := metaPath(s.dir, name)
	m := Meta{Description: description, UpdatedAt: time.Now()}
	if existing, err := s.GetMeta(name); err == nil {
		m.CreatedAt = existing.CreatedAt
	} else {
		m.CreatedAt = m.UpdatedAt
	}
	data, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(p, data, 0600)
}

func (s *Store) GetMeta(name string) (Meta, error) {
	p := metaPath(s.dir, name)
	data, err := os.ReadFile(p)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return Meta{}, nil
		}
		return Meta{}, err
	}
	var m Meta
	if err := json.Unmarshal(data, &m); err != nil {
		return Meta{}, err
	}
	return m, nil
}

func (s *Store) ClearMeta(name string) error {
	p := metaPath(s.dir, name)
	err := os.Remove(p)
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	return err
}
