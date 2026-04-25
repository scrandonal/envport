package store

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
)

func checksumPath(root, name string) string {
	return filepath.Join(root, name+".checksum")
}

// ComputeChecksum computes a deterministic SHA-256 checksum of a snapshot's
// key=value pairs, sorted by key for consistency.
func ComputeChecksum(root, name string) (string, error) {
	snap, err := loadSnapshot(root, name)
	if err != nil {
		return "", err
	}

	keys := make([]string, 0, len(snap.Vars))
	for k := range snap.Vars {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	h := sha256.New()
	for _, k := range keys {
		fmt.Fprintf(h, "%s=%s\n", k, snap.Vars[k])
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

// SaveChecksum computes and persists the checksum for the named snapshot.
func SaveChecksum(root, name string) (string, error) {
	sum, err := ComputeChecksum(root, name)
	if err != nil {
		return "", err
	}

	data, err := json.Marshal(sum)
	if err != nil {
		return "", err
	}

	if err := os.WriteFile(checksumPath(root, name), data, 0600); err != nil {
		return "", err
	}
	return sum, nil
}

// LoadChecksum reads the previously saved checksum for the named snapshot.
func LoadChecksum(root, name string) (string, error) {
	data, err := os.ReadFile(checksumPath(root, name))
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return "", nil
		}
		return "", err
	}

	var sum string
	if err := json.Unmarshal(data, &sum); err != nil {
		return "", err
	}
	return sum, nil
}

// ClearChecksum removes the persisted checksum for the named snapshot.
func ClearChecksum(root, name string) error {
	err := os.Remove(checksumPath(root, name))
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	return err
}

// VerifyChecksum returns true when the snapshot's current content matches its
// stored checksum. Returns false (and no error) when no checksum is stored.
func VerifyChecksum(root, name string) (bool, error) {
	stored, err := LoadChecksum(root, name)
	if err != nil {
		return false, err
	}
	if stored == "" {
		return false, nil
	}
	current, err := ComputeChecksum(root, name)
	if err != nil {
		return false, err
	}
	return current == stored, nil
}
