package store

import (
	"os"
	"path/filepath"
	"testing"
)

func newProjectStore(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	return dir
}

func touchSnapshotForProject(t *testing.T, base, name string) {
	t.Helper()
	f := filepath.Join(base, name+".json")
	if err := os.WriteFile(f, []byte(`{}`), 0600); err != nil {
		t.Fatal(err)
	}
}

func TestSetAndGetProject(t *testing.T) {
	base := newProjectStore(t)
	touchSnapshotForProject(t, base, "mysnap")

	if err := SetProject(base, "mysnap", "alpha"); err != nil {
		t.Fatalf("SetProject: %v", err)
	}
	got, err := GetProject(base, "mysnap")
	if err != nil {
		t.Fatalf("GetProject: %v", err)
	}
	if got != "alpha" {
		t.Errorf("expected %q, got %q", "alpha", got)
	}
}

func TestGetProjectMissing(t *testing.T) {
	base := newProjectStore(t)
	touchSnapshotForProject(t, base, "mysnap")

	got, err := GetProject(base, "mysnap")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "" {
		t.Errorf("expected empty string, got %q", got)
	}
}

func TestClearProject(t *testing.T) {
	base := newProjectStore(t)
	touchSnapshotForProject(t, base, "mysnap")

	_ = SetProject(base, "mysnap", "beta")
	if err := ClearProject(base, "mysnap"); err != nil {
		t.Fatalf("ClearProject: %v", err)
	}
	got, _ := GetProject(base, "mysnap")
	if got != "" {
		t.Errorf("expected empty after clear, got %q", got)
	}
}

func TestClearProjectIdempotent(t *testing.T) {
	base := newProjectStore(t)
	touchSnapshotForProject(t, base, "mysnap")

	if err := ClearProject(base, "mysnap"); err != nil {
		t.Fatalf("ClearProject idempotent: %v", err)
	}
}

func TestSetProjectNotFound(t *testing.T) {
	base := newProjectStore(t)
	err := SetProject(base, "ghost", "alpha")
	if err == nil {
		t.Fatal("expected error for missing snapshot")
	}
}

func TestListByProject(t *testing.T) {
	base := newProjectStore(t)
	for _, name := range []string{"a", "b", "c"} {
		touchSnapshotForProject(t, base, name)
	}
	_ = SetProject(base, "a", "proj-x")
	_ = SetProject(base, "b", "proj-x")
	_ = SetProject(base, "c", "proj-y")

	results, err := ListByProject(base, "proj-x")
	if err != nil {
		t.Fatalf("ListByProject: %v", err)
	}
	if len(results) != 2 {
		t.Errorf("expected 2 results, got %d: %v", len(results), results)
	}
}
