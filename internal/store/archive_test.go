package store

import (
	"testing"
)

func TestArchiveAndList(t *testing.T) {
	s := newTempStore(t)
	vars := map[string]string{"FOO": "bar", "BAZ": "qux"}
	if err := s.Save("snap1", vars); err != nil {
		t.Fatal(err)
	}
	if err := ArchiveSnapshot(s.Base, "snap1", vars); err != nil {
		t.Fatal(err)
	}
	entries, err := ListArchive(s.Base)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if entries[0].Name != "snap1" {
		t.Errorf("expected name snap1, got %s", entries[0].Name)
	}
	if entries[0].Vars["FOO"] != "bar" {
		t.Errorf("expected FOO=bar")
	}
}

func TestArchiveListEmpty(t *testing.T) {
	s := newTempStore(t)
	entries, err := ListArchive(s.Base)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 0 {
		t.Fatalf("expected 0 entries, got %d", len(entries))
	}
}

func TestArchiveMultiple(t *testing.T) {
	s := newTempStore(t)
	vars := map[string]string{"A": "1"}
	for i := 0; i < 3; i++ {
		if err := ArchiveSnapshot(s.Base, "snap", vars); err != nil {
			t.Fatal(err)
		}
	}
	entries, err := ListArchive(s.Base)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(entries))
	}
}

func TestClearArchive(t *testing.T) {
	s := newTempStore(t)
	vars := map[string]string{"X": "y"}
	if err := ArchiveSnapshot(s.Base, "snap", vars); err != nil {
		t.Fatal(err)
	}
	if err := ClearArchive(s.Base); err != nil {
		t.Fatal(err)
	}
	entries, err := ListArchive(s.Base)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 0 {
		t.Fatalf("expected 0 after clear, got %d", len(entries))
	}
}
