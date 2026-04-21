package store

import (
	"testing"
)

func newOwnerStore(t *testing.T) string {
	t.Helper()
	base := t.TempDir()
	s, err := New(base)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	if err := s.Init(); err != nil {
		t.Fatalf("Init: %v", err)
	}
	return base
}

func TestSetAndGetOwner(t *testing.T) {
	base := newOwnerStore(t)
	touchSnapshot(t, base, "mysnap")

	info := OwnerInfo{User: "alice", Email: "alice@example.com", Team: "platform"}
	if err := SetOwner(base, "mysnap", info); err != nil {
		t.Fatalf("SetOwner: %v", err)
	}

	got, err := GetOwner(base, "mysnap")
	if err != nil {
		t.Fatalf("GetOwner: %v", err)
	}
	if got.User != info.User || got.Email != info.Email || got.Team != info.Team {
		t.Errorf("got %+v, want %+v", got, info)
	}
}

func TestGetOwnerMissing(t *testing.T) {
	base := newOwnerStore(t)
	touchSnapshot(t, base, "mysnap")

	got, err := GetOwner(base, "mysnap")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.User != "" {
		t.Errorf("expected empty OwnerInfo, got %+v", got)
	}
}

func TestSetOwnerEmptyUser(t *testing.T) {
	base := newOwnerStore(t)
	touchSnapshot(t, base, "mysnap")

	err := SetOwner(base, "mysnap", OwnerInfo{User: ""})
	if err == nil {
		t.Fatal("expected error for empty user, got nil")
	}
}

func TestSetOwnerNotFound(t *testing.T) {
	base := newOwnerStore(t)

	err := SetOwner(base, "ghost", OwnerInfo{User: "bob"})
	if err == nil {
		t.Fatal("expected error for missing snapshot, got nil")
	}
}

func TestClearOwner(t *testing.T) {
	base := newOwnerStore(t)
	touchSnapshot(t, base, "mysnap")

	if err := SetOwner(base, "mysnap", OwnerInfo{User: "carol"}); err != nil {
		t.Fatalf("SetOwner: %v", err)
	}
	if err := ClearOwner(base, "mysnap"); err != nil {
		t.Fatalf("ClearOwner: %v", err)
	}

	got, err := GetOwner(base, "mysnap")
	if err != nil {
		t.Fatalf("GetOwner after clear: %v", err)
	}
	if got.User != "" {
		t.Errorf("expected empty after clear, got %+v", got)
	}
}

func TestClearOwnerIdempotent(t *testing.T) {
	base := newOwnerStore(t)
	touchSnapshot(t, base, "mysnap")

	if err := ClearOwner(base, "mysnap"); err != nil {
		t.Fatalf("ClearOwner on missing file should not error: %v", err)
	}
}
