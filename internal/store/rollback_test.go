package store_test

import (
	"testing"

	"github.com/nicholasgasior/envport/internal/snapshot"
	"github.com/nicholasgasior/envport/internal/store"
)

func newRollbackStore(t *testing.T) *store.Manager {
	t.Helper()
	s, err := store.New(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	return store.NewManager(s)
}

func TestRollbackSuccess(t *testing.T) {
	m := newRollbackStore(t)

	v1 := snapshot.New(map[string]string{"KEY": "v1"})
	v2 := snapshot.New(map[string]string{"KEY": "v2"})

	if err := m.Save("env", v1); err != nil {
		t.Fatal(err)
	}
	if err := m.Save("env", v2); err != nil {
		t.Fatal(err)
	}

	if err := m.Rollback("env", 1); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	snap, err := m.Load("env")
	if err != nil {
		t.Fatal(err)
	}
	if snap.Vars["KEY"] != "v1" {
		t.Errorf("expected v1, got %q", snap.Vars["KEY"])
	}
}

func TestRollbackNotEnoughHistory(t *testing.T) {
	m := newRollbackStore(t)

	v1 := snapshot.New(map[string]string{"KEY": "v1"})
	if err := m.Save("env", v1); err != nil {
		t.Fatal(err)
	}

	if err := m.Rollback("env", 5); err == nil {
		t.Error("expected error for insufficient history")
	}
}

func TestRollbackNotFound(t *testing.T) {
	m := newRollbackStore(t)
	if err := m.Rollback("missing", 1); err == nil {
		t.Error("expected error for missing snapshot")
	}
}
