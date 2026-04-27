package store

import (
	"os"
	"path/filepath"
	"testing"
)

func newFormatStore(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	return dir
}

func touchSnapshotForFormat(t *testing.T, root, name string) {
	t.Helper()
	f, err := os.Create(filepath.Join(root, name+".json"))
	if err != nil {
		t.Fatal(err)
	}
	f.Close()
}

func TestSetAndGetFormat(t *testing.T) {
	root := newFormatStore(t)
	touchSnapshotForFormat(t, root, "prod")
	if err := SetFormat(root, "prod", "dotenv"); err != nil {
		t.Fatalf("SetFormat: %v", err)
	}
	got, err := GetFormat(root, "prod")
	if err != nil {
		t.Fatalf("GetFormat: %v", err)
	}
	if got != "dotenv" {
		t.Errorf("expected dotenv, got %q", got)
	}
}

func TestGetFormatMissing(t *testing.T) {
	root := newFormatStore(t)
	got, err := GetFormat(root, "ghost")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "" {
		t.Errorf("expected empty string, got %q", got)
	}
}

func TestSetFormatInvalid(t *testing.T) {
	root := newFormatStore(t)
	touchSnapshotForFormat(t, root, "prod")
	if err := SetFormat(root, "prod", "xml"); err == nil {
		t.Error("expected error for invalid format")
	}
}

func TestSetFormatNotFound(t *testing.T) {
	root := newFormatStore(t)
	if err := SetFormat(root, "missing", "json"); err == nil {
		t.Error("expected error for missing snapshot")
	}
}

func TestClearFormat(t *testing.T) {
	root := newFormatStore(t)
	touchSnapshotForFormat(t, root, "dev")
	_ = SetFormat(root, "dev", "shell")
	if err := ClearFormat(root, "dev"); err != nil {
		t.Fatalf("ClearFormat: %v", err)
	}
	got, _ := GetFormat(root, "dev")
	if got != "" {
		t.Errorf("expected empty after clear, got %q", got)
	}
}

func TestClearFormatIdempotent(t *testing.T) {
	root := newFormatStore(t)
	if err := ClearFormat(root, "none"); err != nil {
		t.Errorf("ClearFormat on missing should not error: %v", err)
	}
}

func TestListByFormat(t *testing.T) {
	root := newFormatStore(t)
	for _, n := range []string{"a", "b", "c"} {
		touchSnapshotForFormat(t, root, n)
	}
	_ = SetFormat(root, "a", "json")
	_ = SetFormat(root, "b", "dotenv")
	_ = SetFormat(root, "c", "json")
	names, err := ListByFormat(root, "json")
	if err != nil {
		t.Fatalf("ListByFormat: %v", err)
	}
	if len(names) != 2 {
		t.Errorf("expected 2, got %d: %v", len(names), names)
	}
}
