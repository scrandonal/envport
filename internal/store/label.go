package store

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

func labelPath(root, name string) string {
	return filepath.Join(root, name+".labels.json")
}

func loadLabels(path string) ([]string, error) {
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return []string{}, nil
	}
	if err != nil {
		return nil, err
	}
	var labels []string
	return labels, json.Unmarshal(data, &labels)
}

func saveLabels(path string, labels []string) error {
	data, err := json.Marshal(labels)
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0600)
}

func (s *Store) AddLabel(name, label string) error {
	if !s.Exists(name) {
		return ErrNotFound
	}
	p := labelPath(s.root, name)
	labels, err := loadLabels(p)
	if err != nil {
		return err
	}
	for _, l := range labels {
		if l == label {
			return nil
		}
	}
	return saveLabels(p, append(labels, label))
}

func (s *Store) RemoveLabel(name, label string) error {
	p := labelPath(s.root, name)
	labels, err := loadLabels(p)
	if err != nil {
		return err
	}
	filtered := labels[:0]
	for _, l := range labels {
		if l != label {
			filtered = append(filtered, l)
		}
	}
	return saveLabels(p, filtered)
}

func (s *Store) GetLabels(name string) ([]string, error) {
	return loadLabels(labelPath(s.root, name))
}

func (s *Store) ListByLabel(label string) ([]string, error) {
	names, err := s.List()
	if err != nil {
		return nil, err
	}
	var matched []string
	for _, n := range names {
		labels, err := s.GetLabels(n)
		if err != nil {
			continue
		}
		for _, l := range labels {
			if l == label {
				matched = append(matched, n)
				break
			}
		}
	}
	return matched, nil
}
