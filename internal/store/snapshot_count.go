package store

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

type AccessCount struct {
	Loads  int `json:"loads"`
	Saves  int `json:"saves"`
}

func countPath(base, name string) string {
	return filepath.Join(base, name+".count.json")
}

func IncrementLoad(base, name string) error {
	return modifyCount(base, name, func(c *AccessCount) { c.Loads++ })
}

func IncrementSave(base, name string) error {
	return modifyCount(base, name, func(c *AccessCount) { c.Saves++ })
}

func GetCount(base, name string) (AccessCount, error) {
	path := countPath(base, name)
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return AccessCount{}, nil
	}
	if err != nil {
		return AccessCount{}, err
	}
	var c AccessCount
	if err := json.Unmarshal(data, &c); err != nil {
		return AccessCount{}, err
	}
	return c, nil
}

func ClearCount(base, name string) error {
	path := countPath(base, name)
	err := os.Remove(path)
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	return err
}

func modifyCount(base, name string, fn func(*AccessCount)) error {
	c, err := GetCount(base, name)
	if err != nil {
		return err
	}
	fn(&c)
	data, err := json.Marshal(c)
	if err != nil {
		return err
	}
	return os.WriteFile(countPath(base, name), data, 0600)
}
