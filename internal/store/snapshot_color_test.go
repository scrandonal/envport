package store

import (
	"os"
	"path/filepath"
	"testing"
)

func newColorStore(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	return dir
}

func touchSnapshotForColor(t *testing.T, base, name string) {
	t.Helper()
	f := filepath.Join(base, name+".json")
	if err := os.WriteFile(f, []byte(`{}`), 0600); err != nil {
		t.Fatal(err)
	}
}

func TestSetAndGetColor(t *testing.T) {
	base := newColorStore(t)
	touchSnapshotForColor(t, base, "dev")

	if err := SetColor(base, "dev", "blue"); err != nil {
		t.Fatalf("SetColor: %v", err)
	}
	c, err := GetColor(base, "dev")
	if err != nil {
		t.Fatalf("GetColor: %v", err)
	}
	if c != "blue" {
		t.Errorf("expected blue, got %q", c)
	}
}

func TestGetColorMissing(t *testing.T) {
	base := newColorStore(t)
	c, err := GetColor(base, "nope")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c != "" {
		t.Errorf("expected empty string, got %q", c)
	}
}

func TestSetColorInvalid(t *testing.T) {
	base := newColorStore(t)
	touchSnapshotForColor(t, base, "dev")
	if err := SetColor(base, "dev", "purple"); err == nil {
		t.Error("expected error for invalid color")
	}
}

func TestSetColorNotFound(t *testing.T) {
	base := newColorStore(t)
	if err := SetColor(base, "ghost", "red"); err == nil {
		t.Error("expected error for missing snapshot")
	}
}

func TestClearColor(t *testing.T) {
	base := newColorStore(t)
	touchSnapshotForColor(t, base, "dev")
	_ = SetColor(base, "dev", "green")
	if err := ClearColor(base, "dev"); err != nil {
		t.Fatalf("ClearColor: %v", err)
	}
	c, _ := GetColor(base, "dev")
	if c != "" {
		t.Errorf("expected empty after clear, got %q", c)
	}
}

func TestClearColorIdempotent(t *testing.T) {
	base := newColorStore(t)
	if err := ClearColor(base, "nonexistent"); err != nil {
		t.Errorf("expected no error on double clear, got %v", err)
	}
}

func TestListByColor(t *testing.T) {
	base := newColorStore(t)
	for _, n := range []string{"a", "b", "c"} {
		touchSnapshotForColor(t, base, n)
	}
	_ = SetColor(base, "a", "red")
	_ = SetColor(base, "b", "red")
	_ = SetColor(base, "c", "blue")

	names, err := ListByColor(base, "red")
	if err != nil {
		t.Fatalf("ListByColor: %v", err)
	}
	if len(names) != 2 {
		t.Errorf("expected 2 red snapshots, got %d: %v", len(names), names)
	}
}

func TestListByColorEmpty(t *testing.T) {
	base := newColorStore(t)
	for _, n := range []string{"x", "y"} {
		touchSnapshotForColor(t, base, n)
	}
	_ = SetColor(base, "x", "blue")
	_ = SetColor(base, "y", "blue")

	names, err := ListByColor(base, "red")
	if err != nil {
		t.Fatalf("ListByColor: %v", err)
	}
	if len(names) != 0 {
		t.Errorf("expected 0 red snapshots, got %d: %v", len(names), names)
	}
}
