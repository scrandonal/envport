package store

import (
	"os"
	"path/filepath"
	"testing"
)

func newDescriptionStore(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	// create a dummy snapshot file so existence checks pass
	name := "mysnap"
	_ = os.WriteFile(filepath.Join(dir, name+".json"), []byte(`{}`), 0600)
	return dir
}

func TestSetAndGetDescription(t *testing.T) {
	dir := newDescriptionStore(t)
	if err := SetDescription(dir, "mysnap", "my test snapshot"); err != nil {
		t.Fatalf("SetDescription: %v", err)
	}
	got, err := GetDescription(dir, "mysnap")
	if err != nil {
		t.Fatalf("GetDescription: %v", err)
	}
	if got != "my test snapshot" {
		t.Errorf("expected 'my test snapshot', got %q", got)
	}
}

func TestGetDescriptionMissing(t *testing.T) {
	dir := newDescriptionStore(t)
	got, err := GetDescription(dir, "mysnap")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "" {
		t.Errorf("expected empty string, got %q", got)
	}
}

func TestClearDescription(t *testing.T) {
	dir := newDescriptionStore(t)
	_ = SetDescription(dir, "mysnap", "to be cleared")
	if err := ClearDescription(dir, "mysnap"); err != nil {
		t.Fatalf("ClearDescription: %v", err)
	}
	got, _ := GetDescription(dir, "mysnap")
	if got != "" {
		t.Errorf("expected empty after clear, got %q", got)
	}
}

func TestClearDescriptionIdempotent(t *testing.T) {
	dir := newDescriptionStore(t)
	if err := ClearDescription(dir, "mysnap"); err != nil {
		t.Fatalf("expected no error on missing file, got %v", err)
	}
}

func TestSetDescriptionNotFound(t *testing.T) {
	dir := t.TempDir()
	err := SetDescription(dir, "ghost", "some text")
	if err != ErrNotFound {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}
