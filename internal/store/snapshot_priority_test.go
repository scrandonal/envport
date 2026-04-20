package store

import (
	"os"
	"path/filepath"
	"testing"
)

func newPriorityStore(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	return dir
}

func touchSnapshotForPriority(t *testing.T, root, name string) {
	t.Helper()
	f, err := os.Create(filepath.Join(root, name+".json"))
	if err != nil {
		t.Fatal(err)
	}
	f.Close()
}

func TestSetAndGetPriority(t *testing.T) {
	root := newPriorityStore(t)
	touchSnapshotForPriority(t, root, "prod")

	if err := SetPriority(root, "prod", PriorityHigh); err != nil {
		t.Fatalf("SetPriority: %v", err)
	}
	level, err := GetPriority(root, "prod")
	if err != nil {
		t.Fatalf("GetPriority: %v", err)
	}
	if level != PriorityHigh {
		t.Errorf("expected %q, got %q", PriorityHigh, level)
	}
}

func TestGetPriorityMissing(t *testing.T) {
	root := newPriorityStore(t)
	level, err := GetPriority(root, "nonexistent")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if level != PriorityNormal {
		t.Errorf("expected default %q, got %q", PriorityNormal, level)
	}
}

func TestSetPriorityInvalid(t *testing.T) {
	root := newPriorityStore(t)
	touchSnapshotForPriority(t, root, "dev")
	if err := SetPriority(root, "dev", "urgent"); err == nil {
		t.Error("expected error for invalid priority level")
	}
}

func TestSetPriorityNotFound(t *testing.T) {
	root := newPriorityStore(t)
	if err := SetPriority(root, "ghost", PriorityLow); err == nil {
		t.Error("expected error for missing snapshot")
	}
}

func TestClearPriority(t *testing.T) {
	root := newPriorityStore(t)
	touchSnapshotForPriority(t, root, "staging")

	_ = SetPriority(root, "staging", PriorityCritical)
	if err := ClearPriority(root, "staging"); err != nil {
		t.Fatalf("ClearPriority: %v", err)
	}
	level, _ := GetPriority(root, "staging")
	if level != PriorityNormal {
		t.Errorf("expected default after clear, got %q", level)
	}
}

func TestClearPriorityIdempotent(t *testing.T) {
	root := newPriorityStore(t)
	if err := ClearPriority(root, "nobody"); err != nil {
		t.Errorf("expected no error on missing file, got %v", err)
	}
}

func TestListByPriority(t *testing.T) {
	root := newPriorityStore(t)
	for _, name := range []string{"a", "b", "c"} {
		touchSnapshotForPriority(t, root, name)
	}
	_ = SetPriority(root, "a", PriorityHigh)
	_ = SetPriority(root, "b", PriorityHigh)
	_ = SetPriority(root, "c", PriorityLow)

	high, err := ListByPriority(root, PriorityHigh)
	if err != nil {
		t.Fatalf("ListByPriority: %v", err)
	}
	if len(high) != 2 {
		t.Errorf("expected 2 high-priority snapshots, got %d", len(high))
	}
	low, _ := ListByPriority(root, PriorityLow)
	if len(low) != 1 {
		t.Errorf("expected 1 low-priority snapshot, got %d", len(low))
	}
}
