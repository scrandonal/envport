package store

import (
	"os"
	"path/filepath"
	"testing"
)

func newLifecycleStore(t *testing.T) string {
	t.Helper()
	return t.TempDir()
}

func touchSnapshotForLifecycle(t *testing.T, root, name string) {
	t.Helper()
	f := filepath.Join(root, name+".json")
	if err := os.WriteFile(f, []byte(`{}`), 0600); err != nil {
		t.Fatalf("touch snapshot: %v", err)
	}
}

func TestSetAndGetLifecycle(t *testing.T) {
	root := newLifecycleStore(t)
	touchSnapshotForLifecycle(t, root, "prod")

	if err := SetLifecycle(root, "prod", "active"); err != nil {
		t.Fatalf("SetLifecycle: %v", err)
	}
	lc, err := GetLifecycle(root, "prod")
	if err != nil {
		t.Fatalf("GetLifecycle: %v", err)
	}
	if lc.Stage != "active" {
		t.Errorf("expected stage 'active', got %q", lc.Stage)
	}
	if lc.UpdatedAt.IsZero() {
		t.Error("expected non-zero UpdatedAt")
	}
}

func TestGetLifecycleMissing(t *testing.T) {
	root := newLifecycleStore(t)
	touchSnapshotForLifecycle(t, root, "snap")

	lc, err := GetLifecycle(root, "snap")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if lc.Stage != "" {
		t.Errorf("expected empty stage, got %q", lc.Stage)
	}
}

func TestSetLifecycleInvalid(t *testing.T) {
	root := newLifecycleStore(t)
	touchSnapshotForLifecycle(t, root, "snap")

	if err := SetLifecycle(root, "snap", "unknown"); err == nil {
		t.Error("expected error for invalid stage")
	}
}

func TestSetLifecycleNotFound(t *testing.T) {
	root := newLifecycleStore(t)
	if err := SetLifecycle(root, "ghost", "active"); err == nil {
		t.Error("expected error for missing snapshot")
	}
}

func TestClearLifecycle(t *testing.T) {
	root := newLifecycleStore(t)
	touchSnapshotForLifecycle(t, root, "snap")
	_ = SetLifecycle(root, "snap", "deprecated")

	if err := ClearLifecycle(root, "snap"); err != nil {
		t.Fatalf("ClearLifecycle: %v", err)
	}
	lc, _ := GetLifecycle(root, "snap")
	if lc.Stage != "" {
		t.Errorf("expected empty stage after clear, got %q", lc.Stage)
	}
}

func TestClearLifecycleIdempotent(t *testing.T) {
	root := newLifecycleStore(t)
	if err := ClearLifecycle(root, "nonexistent"); err != nil {
		t.Errorf("expected no error on idempotent clear, got %v", err)
	}
}

func TestListByLifecycle(t *testing.T) {
	root := newLifecycleStore(t)
	for _, n := range []string{"a", "b", "c"} {
		touchSnapshotForLifecycle(t, root, n)
	}
	_ = SetLifecycle(root, "a", "active")
	_ = SetLifecycle(root, "b", "deprecated")
	_ = SetLifecycle(root, "c", "active")

	names, err := ListByLifecycle(root, "active")
	if err != nil {
		t.Fatalf("ListByLifecycle: %v", err)
	}
	if len(names) != 2 {
		t.Errorf("expected 2 active snapshots, got %d", len(names))
	}
}
