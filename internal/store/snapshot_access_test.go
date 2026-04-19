package store_test

import (
	"testing"

	"github.com/nicholasgasior/envport/internal/store"
)

func newAccessStore(t *testing.T) (string, func()) {
	t.Helper()
	s, cleanup := newTempStore(t)
	if err := s.Init(); err != nil {
		t.Fatal(err)
	}
	return s.Root(), cleanup
}

func TestGetAccessMissing(t *testing.T) {
	root, cleanup := newAccessStore(t)
	defer cleanup()

	rec, err := store.GetAccess(root, "missing")
	if err != nil {
		t.Fatal(err)
	}
	if rec.LoadCount != 0 || rec.SaveCount != 0 {
		t.Error("expected zero counts for missing record")
	}
}

func TestRecordLoadAndSave(t *testing.T) {
	root, cleanup := newAccessStore(t)
	defer cleanup()

	if err := store.RecordLoad(root, "snap"); err != nil {
		t.Fatal(err)
	}
	if err := store.RecordLoad(root, "snap"); err != nil {
		t.Fatal(err)
	}
	if err := store.RecordSave(root, "snap"); err != nil {
		t.Fatal(err)
	}

	rec, err := store.GetAccess(root, "snap")
	if err != nil {
		t.Fatal(err)
	}
	if rec.LoadCount != 2 {
		t.Errorf("expected load count 2, got %d", rec.LoadCount)
	}
	if rec.SaveCount != 1 {
		t.Errorf("expected save count 1, got %d", rec.SaveCount)
	}
	if rec.LastLoaded == nil {
		t.Error("expected LastLoaded to be set")
	}
	if rec.LastSaved == nil {
		t.Error("expected LastSaved to be set")
	}
}

func TestClearAccess(t *testing.T) {
	root, cleanup := newAccessStore(t)
	defer cleanup()

	_ = store.RecordLoad(root, "snap")
	if err := store.ClearAccess(root, "snap"); err != nil {
		t.Fatal(err)
	}
	rec, err := store.GetAccess(root, "snap")
	if err != nil {
		t.Fatal(err)
	}
	if rec.LoadCount != 0 {
		t.Error("expected zero counts after clear")
	}
}

func TestClearAccessIdempotent(t *testing.T) {
	root, cleanup := newAccessStore(t)
	defer cleanup()

	if err := store.ClearAccess(root, "nonexistent"); err != nil {
		t.Errorf("expected no error clearing missing access, got %v", err)
	}
}
