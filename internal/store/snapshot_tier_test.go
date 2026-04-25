package store

import (
	"os"
	"path/filepath"
	"testing"
)

func newTierStore(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	return dir
}

func touchSnapshotForTier(t *testing.T, root, name string) {
	t.Helper()
	f := filepath.Join(root, name+".json")
	if err := os.WriteFile(f, []byte(`{}`), 0600); err != nil {
		t.Fatal(err)
	}
}

func TestSetAndGetTier(t *testing.T) {
	root := newTierStore(t)
	touchSnapshotForTier(t, root, "prod")
	if err := SetTier(root, "prod", "premium"); err != nil {
		t.Fatalf("SetTier: %v", err)
	}
	got, err := GetTier(root, "prod")
	if err != nil {
		t.Fatalf("GetTier: %v", err)
	}
	if got != "premium" {
		t.Errorf("expected premium, got %q", got)
	}
}

func TestGetTierMissing(t *testing.T) {
	root := newTierStore(t)
	touchSnapshotForTier(t, root, "prod")
	got, err := GetTier(root, "prod")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "" {
		t.Errorf("expected empty, got %q", got)
	}
}

func TestSetTierInvalid(t *testing.T) {
	root := newTierStore(t)
	touchSnapshotForTier(t, root, "prod")
	if err := SetTier(root, "prod", "gold"); err == nil {
		t.Error("expected error for invalid tier")
	}
}

func TestSetTierNotFound(t *testing.T) {
	root := newTierStore(t)
	if err := SetTier(root, "missing", "free"); err == nil {
		t.Error("expected error for missing snapshot")
	}
}

func TestClearTier(t *testing.T) {
	root := newTierStore(t)
	touchSnapshotForTier(t, root, "prod")
	_ = SetTier(root, "prod", "standard")
	if err := ClearTier(root, "prod"); err != nil {
		t.Fatalf("ClearTier: %v", err)
	}
	got, _ := GetTier(root, "prod")
	if got != "" {
		t.Errorf("expected empty after clear, got %q", got)
	}
}

func TestClearTierIdempotent(t *testing.T) {
	root := newTierStore(t)
	touchSnapshotForTier(t, root, "prod")
	if err := ClearTier(root, "prod"); err != nil {
		t.Errorf("expected no error on double clear: %v", err)
	}
}

func TestListByTier(t *testing.T) {
	root := newTierStore(t)
	for _, n := range []string{"a", "b", "c"} {
		touchSnapshotForTier(t, root, n)
	}
	_ = SetTier(root, "a", "premium")
	_ = SetTier(root, "b", "free")
	_ = SetTier(root, "c", "premium")
	names, err := ListByTier(root, "premium")
	if err != nil {
		t.Fatalf("ListByTier: %v", err)
	}
	if len(names) != 2 {
		t.Errorf("expected 2 premium snapshots, got %d", len(names))
	}
}
