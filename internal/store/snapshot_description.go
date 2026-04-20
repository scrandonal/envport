package store

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

type Description struct {
	Text string `json:"text"`
}

func descriptionPath(base, name string) string {
	return filepath.Join(base, name+".description.json")
}

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

func ClearDescription(base, name string) error {
	err := os.Remove(descriptionPath(base, name))
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	return err
}
