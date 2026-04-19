package store

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

type SnapshotSize struct {
	VarCount  int   `json:"var_count"`
	ByteSize  int64 `json:"byte_size"`
}

func sizePath(base, name string) string {
	return filepath.Join(base, name+".size.json")
}

func ComputeSize(base, name string) (*SnapshotSize, error) {
	snap, err := loadSnapshot(base, name)
	if err != nil {
		return nil, err
	}

	data, err := json.Marshal(snap.Vars)
	if err != nil {
		return nil, err
	}

	return &SnapshotSize{
		VarCount: len(snap.Vars),
		ByteSize: int64(len(data)),
	}, nil
}

func SaveSize(base, name string, sz *SnapshotSize) error {
	data, err := json.MarshalIndent(sz, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(sizePath(base, name), data, 0600)
}

func LoadSize(base, name string) (*SnapshotSize, error) {
	data, err := os.ReadFile(sizePath(base, name))
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return &SnapshotSize{}, nil
		}
		return nil, err
	}
	var sz SnapshotSize
	if err := json.Unmarshal(data, &sz); err != nil {
		return nil, err
	}
	return &sz, nil
}

func ClearSize(base, name string) error {
	err := os.Remove(sizePath(base, name))
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	return err
}
