package store

import (
	"os"
	"path/filepath"
	"testing"
)

func newPlatformStore(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	return dir
}

func touchSnapshotForPlatform(t *testing.T, base, name string) {
	t.Helper()
	f := filepath.Join(base, name+".json")
	if err := os.WriteFile(f, []byte(`{}`), 0600); err != nil {
		t.Fatal(err)
	}
}

func TestSetAndGetPlatform(t *testing.T) {
	dir := newPlatformStore(t)
	touchSnapshotForPlatform(t, dir, "snap1")
	if err := SetPlatform(dir, "snap1", "linux"); err != nil {
		t.Fatalf("SetPlatform: %v", err)
	}
	p, err := GetPlatform(dir, "snap1")
	if err != nil {
		t.Fatalf("GetPlatform: %v", err)
	}
	if p != "linux" {
		t.Errorf("expected linux, got %q", p)
	}
}

func TestGetPlatformMissing(t *testing.T) {
	dir := newPlatformStore(t)
	p, err := GetPlatform(dir, "nope")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p != "" {
		t.Errorf("expected empty, got %q", p)
	}
}

func TestSetPlatformInvalid(t *testing.T) {
	dir := newPlatformStore(t)
	touchSnapshotForPlatform(t, dir, "snap1")
	if err := SetPlatform(dir, "snap1", "amiga"); err == nil {
		t.Error("expected error for invalid platform")
	}
}

func TestSetPlatformNotFound(t *testing.T) {
	dir := newPlatformStore(t)
	if err := SetPlatform(dir, "ghost", "linux"); err == nil {
		t.Error("expected error for missing snapshot")
	}
}

func TestClearPlatform(t *testing.T) {
	dir := newPlatformStore(t)
	touchSnapshotForPlatform(t, dir, "snap1")
	_ = SetPlatform(dir, "snap1", "darwin")
	if err := ClearPlatform(dir, "snap1"); err != nil {
		t.Fatalf("ClearPlatform: %v", err)
	}
	p, _ := GetPlatform(dir, "snap1")
	if p != "" {
		t.Errorf("expected empty after clear, got %q", p)
	}
}

func TestClearPlatformIdempotent(t *testing.T) {
	dir := newPlatformStore(t)
	if err := ClearPlatform(dir, "nonexistent"); err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestListByPlatform(t *testing.T) {
	dir := newPlatformStore(t)
	for _, name := range []string{"a", "b", "c"} {
		touchSnapshotForPlatform(t, dir, name)
	}
	_ = SetPlatform(dir, "a", "linux")
	_ = SetPlatform(dir, "b", "darwin")
	_ = SetPlatform(dir, "c", "linux")
	results, err := ListByPlatform(dir, "linux")
	if err != nil {
		t.Fatalf("ListByPlatform: %v", err)
	}
	if len(results) != 2 {
		t.Errorf("expected 2, got %d: %v", len(results), results)
	}
}
