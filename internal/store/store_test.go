package store_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/envport/internal/store"
)

func newTempStore(t *testing.T) *store.Store {
	t.Helper()
	dir := t.TempDir()
	s := store.New(dir)
	if err := s.Init(); err != nil {
		t.Fatalf("Init: %v", err)
	}
	return s
}

func TestInit(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "nested", "envport")
	s := store.New(dir)
	if err := s.Init(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	info, err := os.Stat(dir)
	if err != nil {
		t.Fatalf("dir not created: %v", err)
	}
	if info.Mode().Perm() != 0700 {
		t.Errorf("expected perm 0700, got %v", info.Mode().Perm())
	}
}

func TestListEmpty(t *testing.T) {
	s := newTempStore(t)
	names, err := s.List()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(names) != 0 {
		t.Errorf("expected empty list, got %v", names)
	}
}

func TestListAndDelete(t *testing.T) {
	s := newTempStore(t)

	for _, name := range []string{"dev", "prod", "staging"} {
		if err := os.WriteFile(s.Path(name), []byte(`{"name":"`+name+`"}`), 0600); err != nil {
			t.Fatalf("setup: %v", err)
		}
	}

	names, err := s.List()
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(names) != 3 {
		t.Errorf("expected 3 snapshots, got %d", len(names))
	}

	if err := s.Delete("prod"); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	names, _ = s.List()
	if len(names) != 2 {
		t.Errorf("expected 2 after delete, got %d", len(names))
	}
}

func TestDeleteNotFound(t *testing.T) {
	s := newTempStore(t)
	if err := s.Delete("ghost"); err != store.ErrNotFound {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestExists(t *testing.T) {
	s := newTempStore(t)
	if s.Exists("missing") {
		t.Error("expected false for missing snapshot")
	}
	_ = os.WriteFile(s.Path("present"), []byte(`{}`), 0600)
	if !s.Exists("present") {
		t.Error("expected true for present snapshot")
	}
}
