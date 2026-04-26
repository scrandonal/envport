package store

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

// AccessCount tracks how many times a snapshot has been loaded or saved.
type AccessCount struct {
	Loads  int `json:"loads"`
	Saves  int `json:"saves"`
}

// Total returns the combined number of loads and saves.
func (c AccessCount) Total() int {
	return c.Loads + c.Saves
}

func countPath(base, name string) string {
	return filepath.Join(base, name+".count.json")
}

// IncrementLoad increments the load counter for the named snapshot.
func IncrementLoad(base, name string) error {
	return modifyCount(base, name, func(c *AccessCount) { c.Loads++ })
}

// IncrementSave increments the save counter for the named snapshot.
func IncrementSave(base, name string) error {
	return modifyCount(base, name, func(c *AccessCount) { c.Saves++ })
}

// GetCount returns the current AccessCount for the named snapshot.
// If no count file exists, an empty AccessCount is returned without error.
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

// ClearCount removes the count file for the named snapshot.
// If no count file exists, ClearCount returns nil.
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
