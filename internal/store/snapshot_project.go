package store

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

func projectPath(base, name string) string {
	return filepath.Join(base, name+".project.json")
}

func SetProject(base, name, project string) error {
	if _, err := os.Stat(filepath.Join(base, name+".json")); errors.Is(err, os.ErrNotExist) {
		return ErrNotFound
	}
	data, err := json.Marshal(project)
	if err != nil {
		return err
	}
	return os.WriteFile(projectPath(base, name), data, 0600)
}

func GetProject(base, name string) (string, error) {
	data, err := os.ReadFile(projectPath(base, name))
	if errors.Is(err, os.ErrNotExist) {
		return "", nil
	}
	if err != nil {
		return "", err
	}
	var project string
	if err := json.Unmarshal(data, &project); err != nil {
		return "", err
	}
	return project, nil
}

func ClearProject(base, name string) error {
	err := os.Remove(projectPath(base, name))
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	return err
}

func ListByProject(base, project string) ([]string, error) {
	entries, err := os.ReadDir(base)
	if err != nil {
		return nil, err
	}
	var results []string
	for _, e := range entries {
		if filepath.Ext(e.Name()) != ".json" {
			continue
		}
		ext := filepath.Ext(e.Name())
		if ext != ".json" {
			continue
		}
		base2 := e.Name()[:len(e.Name())-len(".json")]
		if filepath.Ext(base2) != "" {
			continue
		}
		p, err := GetProject(base, base2)
		if err != nil {
			continue
		}
		if p == project {
			results = append(results, base2)
		}
	}
	return results, nil
}
