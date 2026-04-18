package store

import (
	"testing"
	"time"
)

func TestSetAndGetExpiry(t *testing.T) {
	s := newTempStore(t)
	saveSnap(t, s, "mysnap")

	if err := s.SetExpiry("mysnap", 10*time.Minute); err != nil {
		t.Fatalf("SetExpiry: %v", err)
	}
	e, err := s.GetExpiry("mysnap")
	if err != nil {
		t.Fatalf("GetExpiry: %v", err)
	}
	if e == nil {
		t.Fatal("expected expiry, got nil")
	}
	if time.Until(e.ExpiresAt) <= 0 {
		t.Error("expiry should be in the future")
	}
}

func TestGetExpiryMissing(t *testing.T) {
	s := newTempStore(t)
	saveSnap(t, s, "mysnap")

	e, err := s.GetExpiry("mysnap")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if e != nil {
		t.Error("expected nil expiry for snap with no expiry set")
	}
}

func TestClearExpiry(t *testing.T) {
	s := newTempStore(t)
	saveSnap(t, s, "mysnap")
	_ = s.SetExpiry("mysnap", time.Minute)
	if err := s.ClearExpiry("mysnap"); err != nil {
		t.Fatalf("ClearExpiry: %v", err)
	}
	e, _ := s.GetExpiry("mysnap")
	if e != nil {
		t.Error("expected nil after clear")
	}
}

func TestPruneExpired(t *testing.T) {
	s := newTempStore(t)
	saveSnap(t, s, "old")
	saveSnap(t, s, "fresh")

	_ = s.SetExpiry("old", -time.Second) // already expired
	_ = s.SetExpiry("fresh", time.Hour)

	pruned, err := s.PruneExpired()
	if err != nil {
		t.Fatalf("PruneExpired: %v", err)
	}
	if len(pruned) != 1 || pruned[0] != "old" {
		t.Errorf("expected [old], got %v", pruned)
	}
	if !s.Exists("fresh") {
		t.Error("fresh should still exist")
	}
}

func TestSetExpiryNotFound(t *testing.T) {
	s := newTempStore(t)
	err := s.SetExpiry("ghost", time.Minute)
	if err != ErrNotFound {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func saveSnap(t *testing.T, s *Store, name string) {
	t.Helper()
	if err := s.Save(name, map[string]string{"K": "V"}); err != nil {
		t.Fatalf("save %s: %v", name, err)
	}
}
