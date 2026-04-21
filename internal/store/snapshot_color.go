package store

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

var validColors = map[string]bool{
	"red":    true,
	"green":  true,
	"blue":   true,
	"yellow": true,
	"cyan":   true,
	"magenta": true,
	"white":  true,
	"gray":   true,
}

type colorEntry struct {
	Color string `json:"color"`
}

func colorPath(base, name string) string {
	return filepath.Join(base, name+".color.json")
}

func SetColor(base, name, color string) error {
	if !validColors[color] {
		return fmt.Errorf("invalid color %q: must be one of red, green, blue, yellow, cyan, magenta, white, gray", color)
	}
	if _, err := os.Stat(filepath.Join(base, name+".json")); errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("snapshot %q not found", name)
	}
	data, err := json.Marshal(colorEntry{Color: color})
	if err != nil {
		return err
	}
	return os.WriteFile(colorPath(base, name), data, 0600)
}

func GetColor(base, name string) (string, error) {
	data, err := os.ReadFile(colorPath(base, name))
	if errors.Is(err, os.ErrNotExist) {
		return "", nil
	}
	if err != nil {
		return "", err
	}
	var entry colorEntry
	if err := json.Unmarshal(data, &entry); err != nil {
		return "", err
	}
	return entry.Color, nil
}

func ClearColor(base, name string) error {
	err := os.Remove(colorPath(base, name))
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	return err
}

func ListByColor(base, color string) ([]string, error) {
	entries, err := os.ReadDir(base)
	if err != nil {
		return nil, err
	}
	var names []string
	for _, e := range entries {
		if filepath.Ext(e.Name()) != ".json" {
			continue
		}
		name := e.Name()[:len(e.Name())-len(".json")]
		// skip meta files
		if len(name) > 0 && name[len(name)-1] == ')' {
			continue
		}
		c, err := GetColor(base, name)
		if err != nil || c != color {
			continue
		}
		names = append(names, name)
	}
	return names, nil
}
