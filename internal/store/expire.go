package store

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"time"
)

var ErrAlreadyExpired = errors.New("snapshot has already expired")

type Expiry struct {
	ExpiresAt time.Time `json:"expires_at"`
}

func expirePath(base, name string) string {
	return filepath.Join(base, name+".expire")
}

func (s *Store) SetExpiry(name string, d time.Duration) error {
	if !s.Exists(name) {
		return ErrNotFound
	}
	e := Expiry{ExpiresAt: time.Now().Add(d)}
	data, err := json.Marshal(e)
	if err != nil {
		return err
	}
	return os.WriteFile(expirePath(s.base, name), data, 0600)
}

func (s *Store) GetExpiry(name string) (*Expiry, error) {
	data, err := os.ReadFile(expirePath(s.base, name))
	if errors.Is(err, os.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var e Expiry
	if err := json.Unmarshal(data, &e); err != nil {
		return nil, err
	}
	return &e, nil
}

func (s *Store) ClearExpiry(name string) error {
	err := os.Remove(expirePath(s.base, name))
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	return err
}

func (s *Store) PruneExpired() ([]string, error) {
	names, err := s.List()
	if err != nil {
		return nil, err
	}
	var pruned []string
	for _, name := range names {
		e, err := s.GetExpiry(name)
		if err != nil || e == nil {
			continue
		}
		if time.Now().After(e.ExpiresAt) {
			if err := s.Delete(name); err == nil {
				_ = s.ClearExpiry(name)
				pruned = append(pruned, name)
			}
		}
	}
	return pruned, nil
}
