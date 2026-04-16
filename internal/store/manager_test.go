package store_test

import (
	"os"
	"testing"

	"github.com/user/envport/internal/store"
)

func newManager(t *testing.T) *store.Manager {
	t.Helper()
	s := newTempStore(t)
	return store.NewManager(s)
}

func TestManagerSaveAndLoad(t *testing.T) {
	m := newManager(t)
	env := []string{"FOO=bar", "BAZ=qux"}

	if err := m.Save("test", env); err != nil {
		t.Fatalf("Save: %v", err)
	}

	snap, err := m.Load("test")
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if snap.Name != "test" {
		t.Errorf("expected name %q, got %q", "test", snap.Name)
	}
	if snap.Vars["FOO"] != "bar" {
		t.Errorf("expected FOO=bar, got %q", snap.Vars["FOO"])
	}
}

func TestManagerLoadNotFound(t *testing.T) {
	m := newManager(t)
	_, err := m.Load("ghost")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestManagerList(t *testing.T) {
	m := newManager(t)
	_ = m.Save("alpha", []string{"A=1"})
	_ = m.Save("beta", []string{"B=2"})

	names, err := m.List()
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(names) != 2 {
		t.Errorf("expected 2, got %d", len(names))
	}
}

func TestManagerRename(t *testing.T) {
	m := newManager(t)
	_ = m.Save("old", []string{"X=1"})

	if err := m.Rename("old", "new"); err != nil {
		t.Fatalf("Rename: %v", err)
	}

	snap, err := m.Load("new")
	if err != nil {
		t.Fatalf("Load after rename: %v", err)
	}
	if snap.Name != "new" {
		t.Errorf("expected name %q, got %q", "new", snap.Name)
	}
	if _, err := os.Stat(snap.Name + ".json"); !os.IsNotExist(err) {
		t.Error("old snapshot file should be removed")
	}
}

func TestManagerRenameConflict(t *testing.T) {
	m := newManager(t)
	_ = m.Save("a", []string{})
	_ = m.Save("b", []string{})

	if err := m.Rename("a", "b"); err == nil {
		t.Error("expected error on rename conflict")
	}
}
