package store

import (
	"os"
	"path/filepath"
	"testing"
)

func newProfileStore(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	return dir
}

func touchSnapshotForProfile(t *testing.T, base, name string) {
	t.Helper()
	f := filepath.Join(base, name+".json")
	if err := os.WriteFile(f, []byte(`{}`), 0600); err != nil {
		t.Fatal(err)
	}
}

func TestSetAndGetProfile(t *testing.T) {
	base := newProfileStore(t)
	touchSnapshotForProfile(t, base, "dev")
	p := Profile{Name: "dev", Description: "Development env", Author: "alice"}
	if err := SetProfile(base, "dev", p); err != nil {
		t.Fatal(err)
	}
	got, err := GetProfile(base, "dev")
	if err != nil {
		t.Fatal(err)
	}
	if got.Name != p.Name || got.Author != p.Author || got.Description != p.Description {
		t.Errorf("got %+v, want %+v", got, p)
	}
}

func TestGetProfileMissing(t *testing.T) {
	base := newProfileStore(t)
	touchSnapshotForProfile(t, base, "dev")
	got, err := GetProfile(base, "dev")
	if err != nil {
		t.Fatal(err)
	}
	if got.Name != "" {
		t.Errorf("expected empty profile, got %+v", got)
	}
}

func TestClearProfile(t *testing.T) {
	base := newProfileStore(t)
	touchSnapshotForProfile(t, base, "dev")
	p := Profile{Name: "dev"}
	if err := SetProfile(base, "dev", p); err != nil {
		t.Fatal(err)
	}
	if err := ClearProfile(base, "dev"); err != nil {
		t.Fatal(err)
	}
	got, err := GetProfile(base, "dev")
	if err != nil {
		t.Fatal(err)
	}
	if got.Name != "" {
		t.Errorf("expected empty after clear, got %+v", got)
	}
}

func TestClearProfileIdempotent(t *testing.T) {
	base := newProfileStore(t)
	if err := ClearProfile(base, "missing"); err != nil {
		t.Errorf("expected no error on missing, got %v", err)
	}
}

func TestSetProfileNotFound(t *testing.T) {
	base := newProfileStore(t)
	p := Profile{Name: "ghost"}
	if err := SetProfile(base, "ghost", p); err == nil {
		t.Error("expected error for missing snapshot")
	}
}
