package store

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const (
	MinRating = 1
	MaxRating = 5
)

type Rating struct {
	Value int    `json:"value"`
	Comment string `json:"comment,omitempty"`
}

func ratingPath(base, name string) string {
	return filepath.Join(base, name+".rating.json")
}

func SetRating(base, name string, value int, comment string) error {
	if !snapshotExists(base, name) {
		return fmt.Errorf("snapshot %q not found", name)
	}
	if value < MinRating || value > MaxRating {
		return fmt.Errorf("rating must be between %d and %d", MinRating, MaxRating)
	}
	r := Rating{Value: value, Comment: comment}
	data, err := json.Marshal(r)
	if err != nil {
		return err
	}
	return os.WriteFile(ratingPath(base, name), data, 0600)
}

func GetRating(base, name string) (*Rating, error) {
	data, err := os.ReadFile(ratingPath(base, name))
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var r Rating
	if err := json.Unmarshal(data, &r); err != nil {
		return nil, err
	}
	return &r, nil
}

func ClearRating(base, name string) error {
	err := os.Remove(ratingPath(base, name))
	if os.IsNotExist(err) {
		return nil
	}
	return err
}

func ListByRating(base string, minValue int) ([]string, error) {
	entries, err := os.ReadDir(base)
	if err != nil {
		return nil, err
	}
	var names []string
	for _, e := range entries {
		if filepath.Ext(e.Name()) != ".json" {
			continue
		}
		base2 := e.Name()[:len(e.Name())-len(".rating.json")]
		if len(base2) == 0 || base2 == e.Name() {
			continue
		}
		r, err := GetRating(base, base2)
		if err != nil || r == nil {
			continue
		}
		if r.Value >= minValue {
			names = append(names, base2)
		}
	}
	return names, nil
}
