package store

import (
	"testing"
)

func TestSetAndGetNote(t *testing.T) {
	s := newTempStore(t)
	if err := s.Init("proj"); err != nil {
		t.Fatal(err)
	}
	if err := s.SetNote("proj", "my note"); err != nil {
		t.Fatal(err)
	}
	note, err := s.GetNote("proj")
	if err != nil {
		t.Fatal(err)
	}
	if note != "my note" {
		t.Errorf("expected 'my note', got %q", note)
	}
}

func TestGetNoteMissing(t *testing.T) {
	s := newTempStore(t)
	if err := s.Init("proj"); err != nil {
		t.Fatal(err)
	}
	_, err := s.GetNote("proj")
	if err != ErrNoteNotFound {
		t.Errorf("expected ErrNoteNotFound, got %v", err)
	}
}

func TestClearNote(t *testing.T) {
	s := newTempStore(t)
	if err := s.Init("proj"); err != nil {
		t.Fatal(err)
	}
	s.SetNote("proj", "hello")
	if err := s.ClearNote("proj"); err != nil {
		t.Fatal(err)
	}
	_, err := s.GetNote("proj")
	if err != ErrNoteNotFound {
		t.Errorf("expected ErrNoteNotFound after clear, got %v", err)
	}
}

func TestClearNoteIdempotent(t *testing.T) {
	s := newTempStore(t)
	if err := s.Init("proj"); err != nil {
		t.Fatal(err)
	}
	if err := s.ClearNote("proj"); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestSetNoteNotFound(t *testing.T) {
	s := newTempStore(t)
	err := s.SetNote("ghost", "text")
	if err != ErrNotFound {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}
