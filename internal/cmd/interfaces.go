package cmd

import "github.com/user/envport/internal/snapshot"

// managerIface abstracts store.Manager for testability.
type managerIface interface {
	Save(name string, snap *snapshot.Snapshot) error
	Load(name string) (*snapshot.Snapshot, error)
	List() ([]string, error)
	Delete(name string) error
	Rename(oldName, newName string) error
}
