package store

import (
	"testing"
)

func newAliasStore(t *testing.T) *Store {
	t.Helper()
	s, err := New(t.TempDir())
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	return s
}

func TestSetAndResolveAlias(t *testing.T) {
	s := newAliasStore(t)
	saveSnapshot(t, s, "prod")

	if err := s.SetAlias("p", "prod", false); err != nil {
		t.Fatalf("SetAlias: %v", err)
	}
	got, err := s.ResolveAlias("p")
	if err != nil {
		t.Fatalf("ResolveAlias: %v", err)
	}
	if got != "prod" {
		t.Errorf("expected prod, got %s", got)
	}
}

func TestSetAliasNotFound(t *testing.T) {
	s := newAliasStore(t)
	if err := s.SetAlias("p", "missing", false); err != ErrNotFound {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestSetAliasExists(t *testing.T) {
	s := newAliasStore(t)
	saveSnapshot(t, s, "prod")
	s.SetAlias("p", "prod", false)
	if err := s.SetAlias("p", "prod", false); err != ErrAliasExists {
		t.Errorf("expected ErrAliasExists, got %v", err)
	}
}

func TestSetAliasOverwrite(t *testing.T) {
	s := newAliasStore(t)
	saveSnapshot(t, s, "prod")
	saveSnapshot(t, s, "staging")
	s.SetAlias("p", "prod", false)
	if err := s.SetAlias("p", "staging", true); err != nil {
		t.Fatalf("overwrite: %v", err)
	}
	got, _ := s.ResolveAlias("p")
	if got != "staging" {
		t.Errorf("expected staging, got %s", got)
	}
}

func TestDeleteAlias(t *testing.T) {
	s := newAliasStore(t)
	saveSnapshot(t, s, "prod")
	s.SetAlias("p", "prod", false)
	if err := s.DeleteAlias("p"); err != nil {
		t.Fatalf("DeleteAlias: %v", err)
	}
	if _, err := s.ResolveAlias("p"); err != ErrAliasNotFound {
		t.Errorf("expected ErrAliasNotFound after delete")
	}
}

func TestListAliases(t *testing.T) {
	s := newAliasStore(t)
	saveSnapshot(t, s, "prod")
	saveSnapshot(t, s, "dev")
	s.SetAlias("p", "prod", false)
	s.SetAlias("d", "dev", false)
	m, err := s.ListAliases()
	if err != nil {
		t.Fatalf("ListAliases: %v", err)
	}
	if len(m) != 2 {
		t.Errorf("expected 2 aliases, got %d", len(m))
	}
}

func saveSnapshot(t *testing.T, s *Store, name string) {
	t.Helper()
	if err := s.Save(name, map[string]string{"K": "V"}); err != nil {
		t.Fatalf("Save %s: %v", name, err)
	}
}
