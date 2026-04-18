package store

import (
	"testing"

	"github.com/nicholasgasior/envport/internal/snapshot"
)

func TestMergeAddsNewKeys(t *testing.T) {
	m := newTempStore(t)

	m.Save("base", &snapshot.Snapshot{Vars: map[string]string{"A": "1", "B": "2"}})
	m.Save("extra", &snapshot.Snapshot{Vars: map[string]string{"C": "3"}})

	res, err := m.Merge("extra", "base", MergeSkip)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Added) != 1 || res.Added[0] != "C" {
		t.Errorf("expected Added=[C], got %v", res.Added)
	}

	loaded, _ := m.Load("base")
	if loaded.Vars["C"] != "3" {
		t.Errorf("expected C=3 in merged snapshot")
	}
}

func TestMergeSkipConflict(t *testing.T) {
	m := newTempStore(t)

	m.Save("base", &snapshot.Snapshot{Vars: map[string]string{"A": "original"}})
	m.Save("src", &snapshot.Snapshot{Vars: map[string]string{"A": "new"}})

	res, err := m.Merge("src", "base", MergeSkip)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Skipped) != 1 {
		t.Errorf("expected 1 skipped key, got %v", res.Skipped)
	}

	loaded, _ := m.Load("base")
	if loaded.Vars["A"] != "original" {
		t.Errorf("expected A=original after skip, got %s", loaded.Vars["A"])
	}
}

func TestMergeOverwriteConflict(t *testing.T) {
	m := newTempStore(t)

	m.Save("base", &snapshot.Snapshot{Vars: map[string]string{"A": "original"}})
	m.Save("src", &snapshot.Snapshot{Vars: map[string]string{"A": "new"}})

	res, err := m.Merge("src", "base", MergeOverwrite)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Overwritten) != 1 {
		t.Errorf("expected 1 overwritten key, got %v", res.Overwritten)
	}

	loaded, _ := m.Load("base")
	if loaded.Vars["A"] != "new" {
		t.Errorf("expected A=new after overwrite, got %s", loaded.Vars["A"])
	}
}

func TestMergeSrcNotFound(t *testing.T) {
	m := newTempStore(t)
	m.Save("base", &snapshot.Snapshot{Vars: map[string]string{}})

	_, err := m.Merge("missing", "base", MergeSkip)
	if err == nil {
		t.Error("expected error for missing src")
	}
}
