package store

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

var validPlatforms = map[string]bool{
	"linux":   true,
	"darwin":  true,
	"windows": true,
	"freebsd": true,
	"any":     true,
}

func platformPath(base, name string) string {
	return filepath.Join(base, name+".platform.json")
}

func SetPlatform(base, name, platform string) error {
	if !validPlatforms[platform] {
		return fmt.Errorf("invalid platform %q: must be one of linux, darwin, windows, freebsd, any", platform)
	}
	if _, err := os.Stat(filepath.Join(base, name+".json")); errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("snapshot %q not found", name)
	}
	data, err := json.Marshal(platform)
	if err != nil {
		return err
	}
	return os.WriteFile(platformPath(base, name), data, 0600)
}

func GetPlatform(base, name string) (string, error) {
	data, err := os.ReadFile(platformPath(base, name))
	if errors.Is(err, os.ErrNotExist) {
		return "", nil
	}
	if err != nil {
		return "", err
	}
	var platform string
	if err := json.Unmarshal(data, &platform); err != nil {
		return "", err
	}
	return platform, nil
}

func ClearPlatform(base, name string) error {
	err := os.Remove(platformPath(base, name))
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	return err
}

func ListByPlatform(base, platform string) ([]string, error) {
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
		p, err := GetPlatform(base, base2)
		if err != nil || p != platform {
			continue
		}
		names = append(names, base2)
	}
	return names, nil
}
