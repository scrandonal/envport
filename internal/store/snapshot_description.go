package store

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

// Description holds the human-readable description for a snapshot.
type Description struct {
	Text string `json:"text"`
}

// descriptionPath returns the file path for the description of a snapshot.
func descriptionPath(base, name string) string {
	return filepath.Join(base, name+".description.json")
}

// SetDescription sets the description text for the named snapshot.
// Returns ErrNotFound if the snapshot does not exist.
func SetDescription(base, name, text string) error {
	if _, err := os.Stat(filepath.Join(base, name+".json")); errors.Is(err, os.ErrNotExist) {
		return ErrNotFound
	}
	d := Description{Text: text}
	data, err := json.Marshal(d)
	if err != nil {
		return err
	}
	return os.WriteFile(descriptionPath(base, name), data, 0600)
}

// GetDescription returns the description text for the named snapshot.
// Returns an empty string (and no error) if no description has been set.
func GetDescription(base, name string) (string, error) {
	data, err := os.ReadFile(descriptionPath(base, name))
	if errors.Is(err, os.ErrNotExist) {
		return "", nil
	}
	if err != nil {
		return "", err
	}
	var d Description
	if err := json.Unmarshal(data, &d); err != nil {
		return "", err
	}
	return d.Text, nil
}

// ClearDescription removes the description for the named snapshot.
// Returns nil if no description file exists.
func ClearDescription(base, name string) error {
	err := os.Remove(descriptionPath(base, name))
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	return err
}
