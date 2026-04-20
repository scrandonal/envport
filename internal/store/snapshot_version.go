package store

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// SnapshotVersion records a named version tag for a snapshot.
type SnapshotVersion struct {
	Tag       string    `json:"tag"`
	CreatedAt time.Time `json:"created_at"`
	Note      string    `json:"note,omitempty"`
}

func versionPath(base, name string) string {
	return filepath.Join(base, name+".versions.json")
}

func loadVersions(path string) ([]SnapshotVersion, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return []SnapshotVersion{}, nil
		}
		return nil, err
	}
	var versions []SnapshotVersion
	if err := json.Unmarshal(data, &versions); err != nil {
		return nil, err
	}
	return versions, nil
}

func saveVersions(path string, versions []SnapshotVersion) error {
	data, err := json.MarshalIndent(versions, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0600)
}

// AddVersion tags the current state of a snapshot with a label.
func (s *Store) AddVersion(name, tag, note string) error {
	if !s.Exists(name) {
		return fmt.Errorf("snapshot %q not found", name)
	}
	p := versionPath(s.dir, name)
	versions, err := loadVersions(p)
	if err != nil {
		return err
	}
	for _, v := range versions {
		if v.Tag == tag {
			return fmt.Errorf("version tag %q already exists for snapshot %q", tag, name)
		}
	}
	versions = append(versions, SnapshotVersion{
		Tag:       tag,
		CreatedAt: time.Now().UTC(),
		Note:      note,
	})
	return saveVersions(p, versions)
}

// ListVersions returns all version tags for a snapshot.
func (s *Store) ListVersions(name string) ([]SnapshotVersion, error) {
	if !s.Exists(name) {
		return nil, fmt.Errorf("snapshot %q not found", name)
	}
	return loadVersions(versionPath(s.dir, name))
}

// RemoveVersion deletes a version tag from a snapshot.
func (s *Store) RemoveVersion(name, tag string) error {
	if !s.Exists(name) {
		return fmt.Errorf("snapshot %q not found", name)
	}
	p := versionPath(s.dir, name)
	versions, err := loadVersions(p)
	if err != nil {
		return err
	}
	updated := versions[:0]
	for _, v := range versions {
		if v.Tag != tag {
			updated = append(updated, v)
		}
	}
	if len(updated) == len(versions) {
		return fmt.Errorf("version tag %q not found for snapshot %q", tag, name)
	}
	return saveVersions(p, updated)
}
