package store

import (
	"testing"
)

func newGroupStore(t *testing.T) *Store {
	t.Helper()
	s, err := New(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	return s
}

func TestCreateAndGetGroup(t *testing.T) {
	s := newGroupStore(t)
	saveSnap(t, s, "dev")
	saveSnap(t, s, "prod")

	if err := s.CreateGroup("mygroup", []string{"dev", "prod"}); err != nil {
		t.Fatal(err)
	}
	names, err := s.GetGroup("mygroup")
	if err != nil {
		t.Fatal(err)
	}
	if len(names) != 2 || names[0] != "dev" || names[1] != "prod" {
		t.Fatalf("unexpected group members: %v", names)
	}
}

func TestCreateGroupSnapshotNotFound(t *testing.T) {
	s := newGroupStore(t)
	err := s.CreateGroup("g", []string{"missing"})
	if err != ErrNotFound {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestCreateGroupExists(t *testing.T) {
	s := newGroupStore(t)
	saveSnap(t, s, "dev")
	s.CreateGroup("g", []string{"dev"})
	err := s.CreateGroup("g", []string{"dev"})
	if err != ErrGroupExists {
		t.Fatalf("expected ErrGroupExists, got %v", err)
	}
}

func TestDeleteGroup(t *testing.T) {
	s := newGroupStore(t)
	saveSnap(t, s, "dev")
	s.CreateGroup("g", []string{"dev"})
	if err := s.DeleteGroup("g"); err != nil {
		t.Fatal(err)
	}
	_, err := s.GetGroup("g")
	if err != ErrGroupNotFound {
		t.Fatalf("expected ErrGroupNotFound, got %v", err)
	}
}

func TestListGroups(t *testing.T) {
	s := newGroupStore(t)
	saveSnap(t, s, "a")
	saveSnap(t, s, "b")
	s.CreateGroup("g1", []string{"a"})
	s.CreateGroup("g2", []string{"b"})
	groups, err := s.ListGroups()
	if err != nil {
		t.Fatal(err)
	}
	if len(groups) != 2 {
		t.Fatalf("expected 2 groups, got %d", len(groups))
	}
}

func TestListGroupsEmpty(t *testing.T) {
	s := newGroupStore(t)
	groups, err := s.ListGroups()
	if err != nil {
		t.Fatal(err)
	}
	if len(groups) != 0 {
		t.Fatalf("expected empty, got %v", groups)
	}
}
