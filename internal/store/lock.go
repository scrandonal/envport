package store

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

const lockFileName = ".lock"
const lockTimeout = 5 * time.Second

func lockPath(dir string) string {
	return filepath.Join(dir, lockFileName)
}

// Lock acquires an exclusive lock on the store directory.
// It returns an unlock function and an error.
func (s *Store) Lock() (func(), error) {
	path := lockPath(s.dir)
	deadline := time.Now().Add(lockTimeout)
	for {
		f, err := os.OpenFile(path, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0600)
		if err == nil {
			_ = f.Close()
			return func() { _ = os.Remove(path) }, nil
		}
		if !os.IsExist(err) {
			return nil, fmt.Errorf("lock: %w", err)
		}
		if time.Now().After(deadline) {
			return nil, fmt.Errorf("lock: timed out waiting for store lock")
		}
		time.Sleep(50 * time.Millisecond)
	}
}

// Locked returns true if the store is currently locked.
func (s *Store) Locked() bool {
	_, err := os.Stat(lockPath(s.dir))
	return err == nil
}
