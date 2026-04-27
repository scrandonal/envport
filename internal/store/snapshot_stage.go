package store

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

var validStages = map[string]bool{
	"dev":        true,
	"staging":    true,
	"production": true,
	"test":       true,
	"local":      true,
}

func stagePath(base, name string) string {
	return filepath.Join(base, name+".stage.json")
}

func SetStage(base, name, stage string) error {
	if !validStages[stage] {
		return fmt.Errorf("invalid stage %q: must be one of dev, staging, production, test, local", stage)
	}
	if _, err := os.Stat(filepath.Join(base, name+".json")); errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("snapshot %q not found", name)
	}
	data, err := json.Marshal(stage)
	if err != nil {
		return err
	}
	return os.WriteFile(stagePath(base, name), data, 0600)
}

func GetStage(base, name string) (string, error) {
	data, err := os.ReadFile(stagePath(base, name))
	if errors.Is(err, os.ErrNotExist) {
		return "", nil
	}
	if err != nil {
		return "", err
	}
	var stage string
	if err := json.Unmarshal(data, &stage); err != nil {
		return "", err
	}
	return stage, nil
}

func ClearStage(base, name string) error {
	err := os.Remove(stagePath(base, name))
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	return err
}

func ListByStage(base, stage string) ([]string, error) {
	entries, err := os.ReadDir(base)
	if err != nil {
		return nil, err
	}
	var names []string
	for _, e := range entries {
		if filepath.Ext(e.Name()) != ".json" {
			continue
		}
		base2 := e.Name()[:len(e.Name())-5]
		if filepath.Ext(base2) != "" {
			continue
		}
		s, err := GetStage(base, base2)
		if err != nil || s != stage {
			continue
		}
		names = append(names, base2)
	}
	return names, nil
}
