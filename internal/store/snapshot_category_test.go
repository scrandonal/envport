package store

import (
	"os"
	"path/filepath"
	"testing"
)

func newCategoryStore(t *testing.T) (string, func()) {
	t.Helper()
	dir, err := os.MkdirTemp("", "category-store-*")
	if err != nil {
		t.Fatalf("mkdirtemp: %v", err)
	}
	return dir, func() { os.RemoveAll(dir) }
}

func touchSnapshot(t *testing.T, base, name string) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(base, name+".json"), []byte(`{}`), 0600); err != nil {
		t.Fatalf("touch snapshot: %v", err)
	}
}

func TestSetAndGetCategory(t *testing.T) {
	dir, cleanup := newCategoryStore(t)
	defer cleanup()
	touchSnapshot(t, dir, "prod")

	if err := SetCategory(dir, "prod", "production"); err != nil {
		t.Fatalf("SetCategory: %v", err)
	}
	cat, err := GetCategory(dir, "prod")
	if err != nil {
		t.Fatalf("GetCategory: %v", err)
	}
	if cat != "production" {
		t.Errorf("expected 'production', got %q", cat)
	}
}

func TestGetCategoryMissing(t *testing.T) {
	dir, cleanup := newCategoryStore(t)
	defer cleanup()
	touchSnapshot(t, dir, "dev")

	cat, err := GetCategory(dir, "dev")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cat != "" {
		t.Errorf("expected empty category, got %q", cat)
	}
}

func TestClearCategory(t *testing.T) {
	dir, cleanup := newCategoryStore(t)
	defer cleanup()
	touchSnapshot(t, dir, "staging")

	_ = SetCategory(dir, "staging", "staging")
	if err := ClearCategory(dir, "staging"); err != nil {
		t.Fatalf("ClearCategory: %v", err)
	}
	cat, _ := GetCategory(dir, "staging")
	if cat != "" {
		t.Errorf("expected empty after clear, got %q", cat)
	}
}

func TestClearCategoryIdempotent(t *testing.T) {
	dir, cleanup := newCategoryStore(t)
	defer cleanup()

	if err := ClearCategory(dir, "ghost"); err != nil {
		t.Errorf("ClearCategory on missing should not error: %v", err)
	}
}

func TestSetCategoryNotFound(t *testing.T) {
	dir, cleanup := newCategoryStore(t)
	defer cleanup()

	if err := SetCategory(dir, "missing", "x"); err != ErrNotFound {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestListByCategory(t *testing.T) {
	dir, cleanup := newCategoryStore(t)
	defer cleanup()

	for _, name := range []string{"a", "b", "c"} {
		touchSnapshot(t, dir, name)
	}
	_ = SetCategory(dir, "a", "work")
	_ = SetCategory(dir, "b", "work")
	_ = SetCategory(dir, "c", "personal")

	names, err := ListByCategory(dir, "work")
	if err != nil {
		t.Fatalf("ListByCategory: %v", err)
	}
	if len(names) != 2 {
		t.Errorf("expected 2 results, got %d: %v", len(names), names)
	}
}
