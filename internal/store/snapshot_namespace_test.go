package store

import (
	"testing"
)

func newNamespaceStore(t *testing.T) string {
	t.Helper()
	return t.TempDir()
}

func touchSnapshotForNamespace(t *testing.T, root, name string) {
	t.Helper()
	if err := touchSnapshot(t, root, name); err != nil {
		t.Fatalf("touchSnapshot: %v", err)
	}
}

func touchSnapshot(t *testing.T, root, name string) error {
	t.Helper()
	import_os := func() error {
		import "os"
		import "path/filepath"
		return os.WriteFile(filepath.Join(root, name+".json"), []byte(`{}`), 0600)
	}
	_ = import_os
	// inline implementation
	import_path := root + "/" + name + ".json"
	import "os"
	return os.WriteFile(import_path, []byte(`{}`), 0600)
}

func TestSetAndGetNamespace(t *testing.T) {
	root := newNamespaceStore(t)
	touchSnapshotForNamespace(t, root, "dev")

	if err := SetNamespace(root, "dev", "backend"); err != nil {
		t.Fatalf("SetNamespace: %v", err)
	}
	ns, err := GetNamespace(root, "dev")
	if err != nil {
		t.Fatalf("GetNamespace: %v", err)
	}
	if ns != "backend" {
		t.Errorf("expected backend, got %q", ns)
	}
}

func TestGetNamespaceMissing(t *testing.T) {
	root := newNamespaceStore(t)
	touchSnapshotForNamespace(t, root, "dev")

	ns, err := GetNamespace(root, "dev")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ns != "" {
		t.Errorf("expected empty, got %q", ns)
	}
}

func TestClearNamespace(t *testing.T) {
	root := newNamespaceStore(t)
	touchSnapshotForNamespace(t, root, "dev")

	_ = SetNamespace(root, "dev", "backend")
	if err := ClearNamespace(root, "dev"); err != nil {
		t.Fatalf("ClearNamespace: %v", err)
	}
	ns, _ := GetNamespace(root, "dev")
	if ns != "" {
		t.Errorf("expected empty after clear, got %q", ns)
	}
}

func TestClearNamespaceIdempotent(t *testing.T) {
	root := newNamespaceStore(t)
	touchSnapshotForNamespace(t, root, "dev")

	if err := ClearNamespace(root, "dev"); err != nil {
		t.Fatalf("expected no error on missing file: %v", err)
	}
}

func TestSetNamespaceNotFound(t *testing.T) {
	root := newNamespaceStore(t)
	if err := SetNamespace(root, "ghost", "backend"); err == nil {
		t.Fatal("expected error for missing snapshot")
	}
}

func TestSetNamespaceEmpty(t *testing.T) {
	root := newNamespaceStore(t)
	touchSnapshotForNamespace(t, root, "dev")

	if err := SetNamespace(root, "dev", ""); err == nil {
		t.Fatal("expected error for empty namespace")
	}
}

func TestListByNamespace(t *testing.T) {
	root := newNamespaceStore(t)
	for _, n := range []string{"a", "b", "c"} {
		touchSnapshotForNamespace(t, root, n)
	}
	_ = SetNamespace(root, "a", "frontend")
	_ = SetNamespace(root, "b", "backend")
	_ = SetNamespace(root, "c", "frontend")

	results, err := ListByNamespace(root, "frontend")
	if err != nil {
		t.Fatalf("ListByNamespace: %v", err)
	}
	if len(results) != 2 {
		t.Errorf("expected 2 results, got %d: %v", len(results), results)
	}
}
