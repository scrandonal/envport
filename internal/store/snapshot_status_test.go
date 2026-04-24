package store

import (
	"os"
	"path/filepath"
	"testing"
)

func newStatusStore(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	return dir
}

func touchSnapshotForStatus(t *testing.T, base, name string) {
	t.Helper()
	_ = os.WriteFile(filepath.Join(base, name+".json"), []byte(`{}`), 0600)
}

func TestSetAndGetStatus(t *testing.T) {
	base := newStatusStore(t)
	touchSnapshotForStatus(t, base, "prod")
	if err := SetStatus(base, "prod", "active"); err != nil {
		t.Fatalf("SetStatus: %v", err)
	}
	s, err := GetStatus(base, "prod")
	if err != nil {
		t.Fatalf("GetStatus: %v", err)
	}
	if s != "active" {
		t.Errorf("expected active, got %q", s)
	}
}

func TestGetStatusMissing(t *testing.T) {
	base := newStatusStore(t)
	s, err := GetStatus(base, "ghost")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s != "" {
		t.Errorf("expected empty, got %q", s)
	}
}

func TestSetStatusInvalid(t *testing.T) {
	base := newStatusStore(t)
	touchSnapshotForStatus(t, base, "prod")
	if err := SetStatus(base, "prod", "unknown"); err == nil {
		t.Fatal("expected error for invalid status")
	}
}

func TestSetStatusNotFound(t *testing.T) {
	base := newStatusStore(t)
	if err := SetStatus(base, "missing", "active"); err == nil {
		t.Fatal("expected error for missing snapshot")
	}
}

func TestClearStatus(t *testing.T) {
	base := newStatusStore(t)
	touchSnapshotForStatus(t, base, "prod")
	_ = SetStatus(base, "prod", "draft")
	if err := ClearStatus(base, "prod"); err != nil {
		t.Fatalf("ClearStatus: %v", err)
	}
	s, _ := GetStatus(base, "prod")
	if s != "" {
		t.Errorf("expected empty after clear, got %q", s)
	}
}

func TestClearStatusIdempotent(t *testing.T) {
	base := newStatusStore(t)
	if err := ClearStatus(base, "none"); err != nil {
		t.Fatalf("ClearStatus on missing should not error: %v", err)
	}
}

func TestListByStatus(t *testing.T) {
	base := newStatusStore(t)
	for _, n := range []string{"a", "b", "c"} {
		touchSnapshotForStatus(t, base, n)
	}
	_ = SetStatus(base, "a", "active")
	_ = SetStatus(base, "b", "deprecated")
	_ = SetStatus(base, "c", "active")
	names, err := ListByStatus(base, "active")
	if err != nil {
		t.Fatalf("ListByStatus: %v", err)
	}
	if len(names) != 2 {
		t.Errorf("expected 2 active snapshots, got %d", len(names))
	}
}
