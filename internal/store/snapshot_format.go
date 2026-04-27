package store

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

var validFormats = map[string]bool{
	"dotenv": true,
	"json":   true,
	"shell":  true,
	"yaml":   true,
}

func formatPath(root, name string) string {
	return filepath.Join(root, name+".format")
}

// SetFormat assigns an export format preference to a snapshot.
func SetFormat(root, name, format string) error {
	if !validFormats[format] {
		return fmt.Errorf("invalid format %q: must be one of dotenv, json, shell, yaml", format)
	}
	if _, err := os.Stat(filepath.Join(root, name+".json")); errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("snapshot %q not found", name)
	}
	data, err := json.Marshal(format)
	if err != nil {
		return err
	}
	return os.WriteFile(formatPath(root, name), data, 0600)
}

// GetFormat returns the preferred export format for a snapshot.
// Returns an empty string if no format has been set.
func GetFormat(root, name string) (string, error) {
	data, err := os.ReadFile(formatPath(root, name))
	if errors.Is(err, os.ErrNotExist) {
		return "", nil
	}
	if err != nil {
		return "", err
	}
	var format string
	if err := json.Unmarshal(data, &format); err != nil {
		return "", err
	}
	return format, nil
}

// ClearFormat removes the format preference for a snapshot.
func ClearFormat(root, name string) error {
	err := os.Remove(formatPath(root, name))
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	return err
}

// ListByFormat returns all snapshot names that have the given format set.
func ListByFormat(root, format string) ([]string, error) {
	entries, err := os.ReadDir(root)
	if err != nil {
		return nil, err
	}
	var names []string
	for _, e := range entries {
		if filepath.Ext(e.Name()) != ".format" {
			continue
		}
		name := e.Name()[:len(e.Name())-len(".format")]
		f, err := GetFormat(root, name)
		if err != nil {
			continue
		}
		if f == format {
			names = append(names, name)
		}
	}
	return names, nil
}
