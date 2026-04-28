package store

import (
	"os"
	"path/filepath"
	"testing"
)

func newScopeStore(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	return dir
}

func touchSnapshotForScope(t *testing.T, root, name string) {
	t.Helper()
	f := filepath.Join(root, name+".json")
	if err := os.WriteFile(f, []byte(`{}`), 0600); err != nil {
		t.Fatal(err)
	}
}

func TestSetAndGetScope(t *testing.T) {
	root := newScopeStore(t)
	touchSnapshotForScope(t, root, "dev")
	if err := SetScope(root, "dev", "local"); err != nil {
		t.Fatalf("SetScope: %v", err)
	}
	scope, err := GetScope(root, "dev")
	if err != nil {
		t.Fatalf("GetScope: %v", err)
	}
	if scope != "local" {
		t.Errorf("expected local, got %q", scope)
	}
}

func TestGetScopeMissing(t *testing.T) {
	root := newScopeStore(t)
	scope, err := GetScope(root, "missing")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if scope != "" {
		t.Errorf("expected empty, got %q", scope)
	}
}

func TestSetScopeInvalid(t *testing.T) {
	root := newScopeStore(t)
	touchSnapshotForScope(t, root, "dev")
	if err := SetScope(root, "dev", "unknown"); err == nil {
		t.Fatal("expected error for invalid scope")
	}
}

func TestSetScopeNotFound(t *testing.T) {
	root := newScopeStore(t)
	if err := SetScope(root, "ghost", "global"); err == nil {
		t.Fatal("expected error for missing snapshot")
	}
}

func TestClearScope(t *testing.T) {
	root := newScopeStore(t)
	touchSnapshotForScope(t, root, "dev")
	_ = SetScope(root, "dev", "session")
	if err := ClearScope(root, "dev"); err != nil {
		t.Fatalf("ClearScope: %v", err)
	}
	scope, _ := GetScope(root, "dev")
	if scope != "" {
		t.Errorf("expected empty after clear, got %q", scope)
	}
}

func TestClearScopeIdempotent(t *testing.T) {
	root := newScopeStore(t)
	if err := ClearScope(root, "nonexistent"); err != nil {
		t.Fatalf("ClearScope idempotent: %v", err)
	}
}

func TestListByScope(t *testing.T) {
	root := newScopeStore(t)
	for _, n := range []string{"a", "b", "c"} {
		touchSnapshotForScope(t, root, n)
	}
	_ = SetScope(root, "a", "global")
	_ = SetScope(root, "b", "global")
	_ = SetScope(root, "c", "local")
	names, err := ListByScope(root, "global")
	if err != nil {
		t.Fatalf("ListByScope: %v", err)
	}
	if len(names) != 2 {
		t.Errorf("expected 2 global snapshots, got %d", len(names))
	}
}
