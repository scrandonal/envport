package store_test

import (
	"testing"

	"github.com/nicholasgasior/envport/internal/store"
)

func newCountStore(t *testing.T) string {
	t.Helper()
	s := newTempStore(t)
	if err := s.Init(); err != nil {
		t.Fatalf("init: %v", err)
	}
	return s.Base()
}

func TestGetCountMissing(t *testing.T) {
	base := newCountStore(t)
	c, err := store.GetCount(base, "nosuchsnap")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c.Loads != 0 || c.Saves != 0 {
		t.Fatalf("expected zero counts, got %+v", c)
	}
}

func TestIncrementLoadAndSave(t *testing.T) {
	base := newCountStore(t)
	name := "mysnap"

	for i := 0; i < 3; i++ {
		if err := store.IncrementLoad(base, name); err != nil {
			t.Fatalf("increment load: %v", err)
		}
	}
	if err := store.IncrementSave(base, name); err != nil {
		t.Fatalf("increment save: %v", err)
	}

	c, err := store.GetCount(base, name)
	if err != nil {
		t.Fatalf("get count: %v", err)
	}
	if c.Loads != 3 {
		t.Errorf("expected 3 loads, got %d", c.Loads)
	}
	if c.Saves != 1 {
		t.Errorf("expected 1 save, got %d", c.Saves)
	}
}

func TestClearCount(t *testing.T) {
	base := newCountStore(t)
	name := "mysnap"

	_ = store.IncrementLoad(base, name)
	if err := store.ClearCount(base, name); err != nil {
		t.Fatalf("clear: %v", err)
	}
	c, err := store.GetCount(base, name)
	if err != nil {
		t.Fatalf("get after clear: %v", err)
	}
	if c.Loads != 0 || c.Saves != 0 {
		t.Fatalf("expected zero after clear, got %+v", c)
	}
}

func TestClearCountIdempotent(t *testing.T) {
	base := newCountStore(t)
	if err := store.ClearCount(base, "ghost"); err != nil {
		t.Fatalf("clear nonexistent should not error: %v", err)
	}
}
