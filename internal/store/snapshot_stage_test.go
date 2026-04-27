package store

import (
	"os"
	"path/filepath"
	"testing"
)

func newStageStore(t *testing.T) string {
	t.Helper()
	return t.TempDir()
}

func touchSnapshotForStage(t *testing.T, base, name string) {
	t.Helper()
	f := filepath.Join(base, name+".json")
	if err := os.WriteFile(f, []byte(`{}`), 0600); err != nil {
		t.Fatal(err)
	}
}

func TestSetAndGetStage(t *testing.T) {
	base := newStageStore(t)
	touchSnapshotForStage(t, base, "mysnap")
	if err := SetStage(base, "mysnap", "staging"); err != nil {
		t.Fatal(err)
	}
	s, err := GetStage(base, "mysnap")
	if err != nil {
		t.Fatal(err)
	}
	if s != "staging" {
		t.Errorf("expected staging, got %q", s)
	}
}

func TestGetStageMissing(t *testing.T) {
	base := newStageStore(t)
	s, err := GetStage(base, "ghost")
	if err != nil {
		t.Fatal(err)
	}
	if s != "" {
		t.Errorf("expected empty, got %q", s)
	}
}

func TestSetStageInvalid(t *testing.T) {
	base := newStageStore(t)
	touchSnapshotForStage(t, base, "mysnap")
	if err := SetStage(base, "mysnap", "unknown"); err == nil {
		t.Error("expected error for invalid stage")
	}
}

func TestSetStageNotFound(t *testing.T) {
	base := newStageStore(t)
	if err := SetStage(base, "missing", "dev"); err == nil {
		t.Error("expected error for missing snapshot")
	}
}

func TestClearStage(t *testing.T) {
	base := newStageStore(t)
	touchSnapshotForStage(t, base, "mysnap")
	_ = SetStage(base, "mysnap", "dev")
	if err := ClearStage(base, "mysnap"); err != nil {
		t.Fatal(err)
	}
	s, _ := GetStage(base, "mysnap")
	if s != "" {
		t.Errorf("expected empty after clear, got %q", s)
	}
}

func TestClearStageIdempotent(t *testing.T) {
	base := newStageStore(t)
	if err := ClearStage(base, "nonexistent"); err != nil {
		t.Errorf("expected no error clearing missing stage, got %v", err)
	}
}

func TestListByStage(t *testing.T) {
	base := newStageStore(t)
	for _, n := range []string{"a", "b", "c"} {
		touchSnapshotForStage(t, base, n)
	}
	_ = SetStage(base, "a", "production")
	_ = SetStage(base, "b", "production")
	_ = SetStage(base, "c", "dev")
	names, err := ListByStage(base, "production")
	if err != nil {
		t.Fatal(err)
	}
	if len(names) != 2 {
		t.Errorf("expected 2 production snapshots, got %d", len(names))
	}
}
