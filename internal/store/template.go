package store

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

var ErrTemplateNotFound = errors.New("template not found")

type Template struct {
	Name string            `json:"name"`
	Keys []string          `json:"keys"`
	Defaults map[string]string `json:"defaults,omitempty"`
}

func templatePath(base, name string) string {
	return filepath.Join(base, "templates", name+".json")
}

func (s *Store) SaveTemplate(t Template) error {
	dir := filepath.Join(s.base, "templates")
	if err := os.MkdirAll(dir, 0700); err != nil {
		return fmt.Errorf("mkdir templates: %w", err)
	}
	f, err := os.OpenFile(templatePath(s.base, t.Name), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer f.Close()
	return json.NewEncoder(f).Encode(t)
}

func (s *Store) LoadTemplate(name string) (Template, error) {
	data, err := os.ReadFile(templatePath(s.base, name))
	if errors.Is(err, os.ErrNotExist) {
		return Template{}, ErrTemplateNotFound
	}
	if err != nil {
		return Template{}, err
	}
	var t Template
	if err := json.Unmarshal(data, &t); err != nil {
		return Template{}, err
	}
	return t, nil
}

func (s *Store) DeleteTemplate(name string) error {
	err := os.Remove(templatePath(s.base, name))
	if errors.Is(err, os.ErrNotExist) {
		return ErrTemplateNotFound
	}
	return err
}

func (s *Store) ListTemplates() ([]string, error) {
	dir := filepath.Join(s.base, "templates")
	entries, err := os.ReadDir(dir)
	if errors.Is(err, os.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
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
