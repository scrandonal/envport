package store

import (
	"os"
	"path/filepath"
	"testing"
)

func newDependencyStore(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	return dir
}

func touchSnapshotForDep(t *testing.T, base, name string) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(base, name+".json"), []byte(`{}`), 0600); err != nil {
		t.Fatal(err)
	}
}

func TestSetAndGetDependencies(t *testing.T) {
	base := newDependencyStore(t)
	touchSnapshotForDep(t, base, "app")
	touchSnapshotForDep(t, base, "db")
	touchSnapshotForDep(t, base, "cache")

	if err := SetDependencies(base, "app", []string{"db", "cache"}); err != nil {
		t.Fatalf("SetDependencies: %v", err)
	}
	deps, err := GetDependencies(base, "app")
	if err != nil {
		t.Fatalf("GetDependencies: %v", err)
	}
	if len(deps) != 2 || deps[0] != "db" || deps[1] != "cache" {
		t.Errorf("unexpected deps: %v", deps)
	}
}

func TestGetDependenciesMissing(t *testing.T) {
	base := newDependencyStore(t)
	touchSnapshotForDep(t, base, "app")
	deps, err := GetDependencies(base, "app")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if deps != nil {
		t.Errorf("expected nil, got %v", deps)
	}
}

func TestClearDependencies(t *testing.T) {
	base := newDependencyStore(t)
	touchSnapshotForDep(t, base, "app")
	_ = SetDependencies(base, "app", []string{"db"})
	if err := ClearDependencies(base, "app"); err != nil {
		t.Fatalf("ClearDependencies: %v", err)
	}
	deps, _ := GetDependencies(base, "app")
	if deps != nil {
		t.Errorf("expected nil after clear, got %v", deps)
	}
}

func TestClearDependenciesIdempotent(t *testing.T) {
	base := newDependencyStore(t)
	touchSnapshotForDep(t, base, "app")
	if err := ClearDependencies(base, "app"); err != nil {
		t.Errorf("ClearDependencies on missing file should not error: %v", err)
	}
}

func TestSetDependenciesNotFound(t *testing.T) {
	base := newDependencyStore(t)
	if err := SetDependencies(base, "ghost", []string{"db"}); err == nil {
		t.Error("expected error for missing snapshot")
	}
}

func TestListDependents(t *testing.T) {
	base := newDependencyStore(t)
	touchSnapshotForDep(t, base, "app")
	touchSnapshotForDep(t, base, "worker")
	touchSnapshotForDep(t, base, "db")

	_ = SetDependencies(base, "app", []string{"db"})
	_ = SetDependencies(base, "worker", []string{"db"})

	dependents, err := ListDependents(base, "db")
	if err != nil {
		t.Fatalf("ListDependents: %v", err)
	}
	if len(dependents) != 2 {
		t.Errorf("expected 2 dependents, got %d: %v", len(dependents), dependents)
	}
}
