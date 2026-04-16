package snapshot_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yourorg/envport/internal/snapshot"
)

func TestFromEnviron(t *testing.T) {
	input := []string{"FOO=bar", "BAZ=qux=extra", "EMPTY="}
	got := snapshot.FromEnviron(input)

	cases := map[string]string{
		"FOO":   "bar",
		"BAZ":   "qux=extra",
		"EMPTY": "",
	}
	for k, want := range cases {
		if got[k] != want {
			t.Errorf("key %q: got %q, want %q", k, got[k], want)
		}
	}
}

func TestSaveAndLoad(t *testing.T) {
	env := map[string]string{"API_KEY": "secret", "PORT": "8080"}
	s := snapshot.New("test-snap", env)

	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "snap.json")

	if err := s.Save(path); err != nil {
		t.Fatalf("Save() error: %v", err)
	}

	loaded, err := snapshot.Load(path)
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}

	if loaded.Name != s.Name {
		t.Errorf("Name: got %q, want %q", loaded.Name, s.Name)
	}
	for k, want := range env {
		if loaded.Env[k] != want {
			t.Errorf("Env[%q]: got %q, want %q", k, loaded.Env[k], want)
		}
	}
}

func TestLoadMissingFile(t *testing.T) {
	_, err := snapshot.Load("/nonexistent/path/snap.json")
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}

func TestSaveRestrictsPermissions(t *testing.T) {
	s := snapshot.New("perm-test", map[string]string{"X": "1"})
	path := filepath.Join(t.TempDir(), "snap.json")
	if err := s.Save(path); err != nil {
		t.Fatal(err)
	}
	info, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	if info.Mode().Perm() != 0600 {
		t.Errorf("expected perm 0600, got %v", info.Mode().Perm())
	}
}
