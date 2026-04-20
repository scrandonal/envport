package store

import (
	"os"
	"path/filepath"
	"testing"
)

func newVisibilityStore(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	return dir
}

func touchSnapshotForVisibility(t *testing.T, base, name string) {
	t.Helper()
	f, err := os.Create(filepath.Join(base, name+".json"))
	if err != nil {
		t.Fatal(err)
	}
	f.Close()
}

func TestSetAndGetVisibility(t *testing.T) {
	base := newVisibilityStore(t)
	touchSnapshotForVisibility(t, base, "prod")

	if err := SetVisibility(base, "prod", VisibilityShared); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	v, err := GetVisibility(base, "prod")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v != VisibilityShared {
		t.Errorf("expected shared, got %q", v)
	}
}

func TestGetVisibilityMissing(t *testing.T) {
	base := newVisibilityStore(t)
	v, err := GetVisibility(base, "ghost")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v != VisibilityPrivate {
		t.Errorf("expected default private, got %q", v)
	}
}

func TestSetVisibilityInvalid(t *testing.T) {
	base := newVisibilityStore(t)
	touchSnapshotForVisibility(t, base, "dev")

	if err := SetVisibility(base, "dev", Visibility("secret")); err == nil {
		t.Error("expected error for invalid visibility")
	}
}

func TestSetVisibilityNotFound(t *testing.T) {
	base := newVisibilityStore(t)
	if err := SetVisibility(base, "missing", VisibilityPublic); err == nil {
		t.Error("expected error for missing snapshot")
	}
}

func TestClearVisibility(t *testing.T) {
	base := newVisibilityStore(t)
	touchSnapshotForVisibility(t, base, "staging")

	_ = SetVisibility(base, "staging", VisibilityPublic)
	if err := ClearVisibility(base, "staging"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	v, _ := GetVisibility(base, "staging")
	if v != VisibilityPrivate {
		t.Errorf("expected private after clear, got %q", v)
	}
}

func TestClearVisibilityIdempotent(t *testing.T) {
	base := newVisibilityStore(t)
	if err := ClearVisibility(base, "nope"); err != nil {
		t.Errorf("expected no error on missing file, got %v", err)
	}
}

func TestListByVisibility(t *testing.T) {
	base := newVisibilityStore(t)
	for _, name := range []string{"a", "b", "c"} {
		touchSnapshotForVisibility(t, base, name)
	}
	_ = SetVisibility(base, "a", VisibilityPublic)
	_ = SetVisibility(base, "b", VisibilityShared)
	// c defaults to private

	public, err := ListByVisibility(base, VisibilityPublic, []string{"a", "b", "c"})
	if err != nil {
		t.Fatal(err)
	}
	if len(public) != 1 || public[0] != "a" {
		t.Errorf("expected [a], got %v", public)
	}

	private, err := ListByVisibility(base, VisibilityPrivate, []string{"a", "b", "c"})
	if err != nil {
		t.Fatal(err)
	}
	if len(private) != 1 || private[0] != "c" {
		t.Errorf("expected [c], got %v", private)
	}
}
