package store

import (
	"testing"

	"github.com/user/envport/internal/snapshot"
)

func TestCloneSuccess(t *testing.T) {
	m := newManager(t)
	snap := snapshot.New(map[string]string{"FOO": "bar"})
	if err := m.Save("src", snap); err != nil {
		t.Fatal(err)
	}
	if err := m.Clone("src", "dst", false); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	loaded, err := m.Load("dst")
	if err != nil {
		t.Fatal(err)
	}
	if loaded.Vars["FOO"] != "bar" {
		t.Errorf("expected FOO=bar, got %q", loaded.Vars["FOO"])
	}
}

func TestCloneSrcNotFound(t *testing.T) {
	m := newManager(t)
	err := m.Clone("missing", "dst", false)
	if err == nil {
		t.Fatal("expected error for missing source")
	}
}

func TestCloneDestExists(t *testing.T) {
	m := newManager(t)
	snap := snapshot.New(map[string]string{"A": "1"})
	m.Save("src", snap)
	m.Save("dst", snapshot.New(map[string]string{"B": "2"}))

	err := m.Clone("src", "dst", false)
	if err == nil {
		t.Fatal("expected error when dest exists without overwrite")
	}
}

func TestCloneOverwrite(t *testing.T) {
	m := newManager(t)
	snap := snapshot.New(map[string]string{"X": "new"})
	m.Save("src", snap)
	m.Save("dst", snapshot.New(map[string]string{"X": "old"}))

	if err := m.Clone("src", "dst", true); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	loaded, _ := m.Load("dst")
	if loaded.Vars["X"] != "new" {
		t.Errorf("expected X=new after overwrite, got %q", loaded.Vars["X"])
	}
}
