package store

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type ArchiveEntry struct {
	Name      string            `json:"name"`
	Vars      map[string]string `json:"vars"`
	ArchivedAt time.Time        `json:"archived_at"`
}

func archivePath(base string) string {
	return filepath.Join(base, "archive")
}

func ArchiveSnapshot(base, name string, vars map[string]string) error {
	if err := os.MkdirAll(archivePath(base), 0700); err != nil {
		return err
	}
	entry := ArchiveEntry{
		Name:      name,
		Vars:      vars,
		ArchivedAt: time.Now().UTC(),
	}
	data, err := json.MarshalIndent(entry, "", "  ")
	if err != nil {
		return err
	}
	filename := fmt.Sprintf("%s_%d.json", name, entry.ArchivedAt.UnixNano())
	return os.WriteFile(filepath.Join(archivePath(base), filename), data, 0600)
}

func ListArchive(base string) ([]ArchiveEntry, error) {
	dir := archivePath(base)
	entries, err := os.ReadDir(dir)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var result []ArchiveEntry
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		data, err := os.ReadFile(filepath.Join(dir, e.Name()))
		if err != nil {
			return nil, err
		}
		var ae ArchiveEntry
		if err := json.Unmarshal(data, &ae); err != nil {
			return nil, err
		}
		result = append(result, ae)
	}
	return result, nil
}

func ClearArchive(base string) error {
	dir := archivePath(base)
	if err := os.RemoveAll(dir); err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}
