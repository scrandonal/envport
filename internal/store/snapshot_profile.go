package store

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

type Profile struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Author      string `json:"author,omitempty"`
}

func profilePath(base, name string) string {
	return filepath.Join(base, name+".profile.json")
}

func SetProfile(base, name string, p Profile) error {
	if _, err := os.Stat(filepath.Join(base, name+".json")); errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("snapshot %q not found", name)
	}
	data, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(profilePath(base, name), data, 0600)
}

func GetProfile(base, name string) (Profile, error) {
	data, err := os.ReadFile(profilePath(base, name))
	if errors.Is(err, os.ErrNotExist) {
		return Profile{}, nil
	}
	if err != nil {
		return Profile{}, err
	}
	var p Profile
	if err := json.Unmarshal(data, &p); err != nil {
		return Profile{}, err
	}
	return p, nil
}

func ClearProfile(base, name string) error {
	err := os.Remove(profilePath(base, name))
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	return err
}
