package store

import (
	"os"
	"path/filepath"
	"testing"
)

func newProjectStore(t *testing.T) string {
	t.Helper()
	dir, err := os.MkdirTemp("", "envport-project-*")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.RemoveAll(dir) })
	return dir
}

func touchSnapshotForProject(t *testing.T, base, name string) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(base, name+".json"), []byte("{}"), 0600); err != nil {
		t.Fatal(err)
	}
}

func TestSetAndGetProject(t *testing.T) {
	base := newProjectStore(t)
	touchSnapshotForProject(t, base, "dev")

	if err := SetProject(base, "dev", "myapp"); err != nil {
		t.Fatalf("SetProject: %v", err)
	}
	p, err := GetProject(base, "dev")
	if err != nil {
		t.Fatalf("GetProject: %v", err)
	}
	if p != "myapp" {
		t.Errorf("expected myapp, got %q", p)
	}
}

func TestGetProjectMissing(t *testing.T) {
	base := newProjectStore(t)
	p, err := GetProject(base, "ghost")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p != "" {
		t.Errorf("expected empty, got %q", p)
	}
}

func TestClearProject(t *testing.T) {
	base := newProjectStore(t)
	touchSnapshotForProject(t, base, "dev")
	_ = SetProject(base, "dev", "myapp")

	if err := ClearProject(base, "dev"); err != nil {
		t.Fatalf("ClearProject: %v", err)
	}
	p, _ := GetProject(base, "dev")
	if p != "" {
		t.Errorf("expected empty after clear, got %q", p)
	}
}

func TestClearProjectIdempotent(t *testing.T) {
	base := newProjectStore(t)
	if err := ClearProject(base, "nope"); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestSetProjectNotFound(t *testing.T) {
	base := newProjectStore(t)
	err := SetProject(base, "missing", "myapp")
	if err != ErrNotFound {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestListByProject(t *testing.T) {
	base := newProjectStore(t)
	for _, name := range []string{"dev", "staging", "prod"} {
		touchSnapshotForProject(t, base, name)
	}
	_ = SetProject(base, "dev", "myapp")
	_ = SetProject(base, "staging", "myapp")
	_ = SetProject(base, "prod", "otherapp")

	results, err := ListByProject(base, "myapp")
	if err != nil {
		t.Fatalf("ListByProject: %v", err)
	}
	if len(results) != 2 {
		t.Errorf("expected 2 results, got %d: %v", len(results), results)
	}
}
