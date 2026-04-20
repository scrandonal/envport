package store

import (
	"testing"
)

func newVersionStore(t *testing.T) *Store {
	t.Helper()
	s, err := New(t.TempDir())
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	return s
}

func TestAddAndListVersions(t *testing.T) {
	s := newVersionStore(t)
	if err := s.Init("snap1"); err != nil {
		t.Fatalf("Init: %v", err)
	}
	if err := s.AddVersion("snap1", "v1.0", "initial release"); err != nil {
		t.Fatalf("AddVersion: %v", err)
	}
	if err := s.AddVersion("snap1", "v1.1", "patch"); err != nil {
		t.Fatalf("AddVersion: %v", err)
	}
	versions, err := s.ListVersions("snap1")
	if err != nil {
		t.Fatalf("ListVersions: %v", err)
	}
	if len(versions) != 2 {
		t.Fatalf("expected 2 versions, got %d", len(versions))
	}
	if versions[0].Tag != "v1.0" || versions[1].Tag != "v1.1" {
		t.Errorf("unexpected tags: %v", versions)
	}
}

func TestAddVersionDuplicate(t *testing.T) {
	s := newVersionStore(t)
	if err := s.Init("snap1"); err != nil {
		t.Fatalf("Init: %v", err)
	}
	if err := s.AddVersion("snap1", "v1.0", ""); err != nil {
		t.Fatalf("AddVersion: %v", err)
	}
	if err := s.AddVersion("snap1", "v1.0", ""); err == nil {
		t.Error("expected error for duplicate tag")
	}
}

func TestAddVersionNotFound(t *testing.T) {
	s := newVersionStore(t)
	if err := s.AddVersion("missing", "v1.0", ""); err == nil {
		t.Error("expected error for missing snapshot")
	}
}

func TestListVersionsEmpty(t *testing.T) {
	s := newVersionStore(t)
	if err := s.Init("snap1"); err != nil {
		t.Fatalf("Init: %v", err)
	}
	versions, err := s.ListVersions("snap1")
	if err != nil {
		t.Fatalf("ListVersions: %v", err)
	}
	if len(versions) != 0 {
		t.Errorf("expected 0 versions, got %d", len(versions))
	}
}

func TestRemoveVersion(t *testing.T) {
	s := newVersionStore(t)
	if err := s.Init("snap1"); err != nil {
		t.Fatalf("Init: %v", err)
	}
	_ = s.AddVersion("snap1", "v1.0", "")
	_ = s.AddVersion("snap1", "v1.1", "")
	if err := s.RemoveVersion("snap1", "v1.0"); err != nil {
		t.Fatalf("RemoveVersion: %v", err)
	}
	versions, _ := s.ListVersions("snap1")
	if len(versions) != 1 || versions[0].Tag != "v1.1" {
		t.Errorf("unexpected versions after remove: %v", versions)
	}
}

func TestRemoveVersionNotFound(t *testing.T) {
	s := newVersionStore(t)
	if err := s.Init("snap1"); err != nil {
		t.Fatalf("Init: %v", err)
	}
	if err := s.RemoveVersion("snap1", "v9.9"); err == nil {
		t.Error("expected error for missing version tag")
	}
}
