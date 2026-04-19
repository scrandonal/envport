package store

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

type Bookmark struct {
	Name  string `json:"name"`
	Label string `json:"label"`
}

func bookmarkPath(base string) string {
	return filepath.Join(base, "bookmarks.json")
}

func loadBookmarks(base string) ([]Bookmark, error) {
	path := bookmarkPath(base)
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return []Bookmark{}, nil
	}
	if err != nil {
		return nil, err
	}
	var bm []Bookmark
	return bm, json.Unmarshal(data, &bm)
}

func saveBookmarks(base string, bm []Bookmark) error {
	data, err := json.MarshalIndent(bm, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(bookmarkPath(base), data, 0600)
}

func AddBookmark(base, name, label string) error {
	if _, err := os.Stat(filepath.Join(base, name+".json")); errors.Is(err, os.ErrNotExist) {
		return ErrNotFound
	}
	bm, err := loadBookmarks(base)
	if err != nil {
		return err
	}
	for _, b := range bm {
		if b.Name == name {
			return errors.New("bookmark already exists for: " + name)
		}
	}
	bm = append(bm, Bookmark{Name: name, Label: label})
	return saveBookmarks(base, bm)
}

func RemoveBookmark(base, name string) error {
	bm, err := loadBookmarks(base)
	if err != nil {
		return err
	}
	next := bm[:0]
	for _, b := range bm {
		if b.Name != name {
			next = append(next, b)
		}
	}
	if len(next) == len(bm) {
		return ErrNotFound
	}
	return saveBookmarks(base, next)
}

func ListBookmarks(base string) ([]Bookmark, error) {
	return loadBookmarks(base)
}
