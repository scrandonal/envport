package store

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

type RuntimeInfo struct {
	OS      string `json:"os"`
	Arch    string `json:"arch"`
	Host    string `json:"host"`
	User    string `json:"user"`
	Shell   string `json:"shell"`
}

func runtimePath(base, name string) string {
	return filepath.Join(base, name+".runtime.json")
}

func SetRuntime(base, name string, info RuntimeInfo) error {
	if _, err := os.Stat(filepath.Join(base, name+".json")); errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("snapshot %q not found", name)
	}
	data, err := json.MarshalIndent(info, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(runtimePath(base, name), data, 0600)
}

func GetRuntime(base, name string) (RuntimeInfo, error) {
	data, err := os.ReadFile(runtimePath(base, name))
	if errors.Is(err, os.ErrNotExist) {
		return RuntimeInfo{}, nil
	}
	if err != nil {
		return RuntimeInfo{}, err
	}
	var info RuntimeInfo
	if err := json.Unmarshal(data, &info); err != nil {
		return RuntimeInfo{}, err
	}
	return info, nil
}

func ClearRuntime(base, name string) error {
	err := os.Remove(runtimePath(base, name))
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	return err
}
