package store

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

type ProjectInfo struct {
	Name string `json:"name"`
	URL  string `json:"url,omitempty"`
}

func projectPath(base, name string) string {
	return filepath.Join(base, name+".project.json")
}

func SetProject(base, name string, info ProjectInfo) error {
	if _, err := os.Stat(filepath.Join(base, name+".json")); errors.Is(err, os.ErrNotExist) {
		return ErrNotFound
	}
	data, err := json.MarshalIndent(info, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(projectPath(base, name), data, 0600)
}

func GetProject(base, name string) (ProjectInfo, error) {
	data, err := os.ReadFile(projectPath(base, name))
	if errors.Is(err, os.ErrNotExist) {
		return ProjectInfo{}, nil
	}
	if err != nil {
		return ProjectInfo{}, err
	}
	var info ProjectInfo
	if err := json.Unmarshal(data, &info); err != nil {
		return ProjectInfo{}, err
	}
	return info, nil
}

func ClearProject(base, name string) error {
	err := os.Remove(projectPath(base, name))
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	return err
}

func ListByProject(base, projectName string) ([]string, error) {
	entries, err := os.ReadDir(base)
	if err != nil {
		return nil, err
	}
	var matches []string
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
		info, err := GetProject(base, base2)
		if err != nil {
			continue
		}
		if info.Name == projectName {
			matches = append(matches, base2)
		}
	}
	return matches, nil
}
