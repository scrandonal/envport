package store

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

func projectPath(base, name string) string {
	return filepath.Join(base, name+".project.json")
}

func SetProject(base, name, project string) error {
	if !snapshotExists(base, name) {
		return fmt.Errorf("snapshot %q not found", name)
	}
	if project == "" {
		return errors.New("project name must not be empty")
	}
	data, err := json.Marshal(project)
	if err != nil {
		return err
	}
	return os.WriteFile(projectPath(base, name), data, 0600)
}

func GetProject(base, name string) (string, error) {
	data, err := os.ReadFile(projectPath(base, name))
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return "", nil
		}
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
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}
	return nil
}

func ListByProject(base, project string) ([]string, error) {
	entries, err := os.ReadDir(base)
	if err != nil {
		return nil, err
	}
	var results []string
	for _, e := range entries {
		if e.IsDir() || filepath.Ext(e.Name()) != ".json" {
			continue
		}
		ext := filepath.Ext(e.Name())
		base2 := e.Name()[:len(e.Name())-len(ext)]
		if filepath.Ext(base2) != "" {
			continue
		}
		snap := base2
		p, err := GetProject(base, snap)
		if err != nil || p != project {
			continue
		}
		results = append(results, snap)
	}
	return results, nil
}
