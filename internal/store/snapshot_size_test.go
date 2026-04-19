package store

import (
	"testing"

	"github.com/nicholasgasior/envport/internal/snapshot"
)

func newSizeStore(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	s := New(dir)
	if err := s.Init(); err != nil {
		t.Fatalf("init: %v", err)
	}
	return dir
}

func TestComputeSize(t *testing.T) {
	dir := newSizeStore(t)
	snap := snapshot.New(map[string]string{"FOO": "bar", "BAZ": "qux"})
	if err := saveSnapshot(dir, "mysnap", snap); err != nil {
		t.Fatalf("save: %v", err)
	}

	sz, err := ComputeSize(dir, "mysnap")
	if err != nil {
		t.Fatalf("compute: %v", err)
	}
	if sz.VarCount != 2 {
		t.Errorf("expected 2 vars, got %d", sz.VarCount)
	}
	if sz.ByteSize == 0 {
		t.Error("expected non-zero byte size")
	}
}

func TestSaveAndLoadSize(t *testing.T) {
	dir := newSizeStore(t)
	snap := snapshot.New(map[string]string{"A": "1"})
	if err := saveSnapshot(dir, "s", snap); err != nil {
		t.Fatalf("save: %v", err)
	}

	sz, _ := ComputeSize(dir, "s")
	if err := SaveSize(dir, "s", sz); err != nil {
		t.Fatalf("save size: %v", err)
	}

	loaded, err := LoadSize(dir, "s")
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if loaded.VarCount != sz.VarCount {
		t.Errorf("var count mismatch: got %d want %d", loaded.VarCount, sz.VarCount)
	}
}

func TestLoadSizeMissing(t *testing.T) {
	dir := newSizeStore(t)
	sz, err := LoadSize(dir, "ghost")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if sz.VarCount != 0 || sz.ByteSize != 0 {
		t.Error("expected zero value for missing size")
	}
}

func TestClearSize(t *testing.T) {
	dir := newSizeStore(t)
	snap := snapshot.New(map[string]string{"X": "y"})
	_ = saveSnapshot(dir, "s", snap)
	sz, _ := ComputeSize(dir, "s")
	_ = SaveSize(dir, "s", sz)

	if err := ClearSize(dir, "s"); err != nil {
		t.Fatalf("clear: %v", err)
	}
	loaded, _ := LoadSize(dir, "s")
	if loaded.VarCount != 0 {
		t.Error("expected zero after clear")
	}
}
