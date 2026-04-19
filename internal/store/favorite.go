package store

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

func favoritePath(root string) string {
	return filepath.Join(root, "favorites.json")
}

func loadFavorites(root string) ([]string, error) {
	data, err := os.ReadFile(favoritePath(root))
	if errors.Is(err, os.ErrNotExist) {
		return []string{}, nil
	}
	if err != nil {
		return nil, err
	}
	var favs []string
	if err := json.Unmarshal(data, &favs); err != nil {
		return nil, err
	}
	return favs, nil
}

func saveFavorites(root string, favs []string) error {
	data, err := json.MarshalIndent(favs, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(favoritePath(root), data, 0600)
}

func AddFavorite(root, name string) error {
	if _, err := os.Stat(snapshotPath(root, name)); errors.Is(err, os.ErrNotExist) {
		return ErrNotFound
	}
	favs, err := loadFavorites(root)
	if err != nil {
		return err
	}
	for _, f := range favs {
		if f == name {
			return nil
		}
	}
	return saveFavorites(root, append(favs, name))
}

func RemoveFavorite(root, name string) error {
	favs, err := loadFavorites(root)
	if err != nil {
		return err
	}
	next := favs[:0]
	for _, f := range favs {
		if f != name {
			next = append(next, f)
		}
	}
	return saveFavorites(root, next)
}

func ListFavorites(root string) ([]string, error) {
	return loadFavorites(root)
}
