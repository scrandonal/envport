package store

import (
	"testing"
)

func newLabelStore(t *testing.T) *Store {
	t.Helper()
	s, err := New(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	return s
}

func TestAddAndGetLabels(t *testing.T) {
	s := newLabelStore(t)
	saveTestSnapshot(t, s, "snap1")
	if err := s.AddLabel("snap1", "prod"); err != nil {
		t.Fatal(err)
	}
	if err := s.AddLabel("snap1", "infra"); err != nil {
		t.Fatal(err)
	}
	labels, err := s.GetLabels("snap1")
	if err != nil {
		t.Fatal(err)
	}
	if len(labels) != 2 {
		t.Fatalf("expected 2 labels, got %d", len(labels))
	}
}

func TestAddLabelDuplicate(t *testing.T) {
	s := newLabelStore(t)
	saveTestSnapshot(t, s, "snap1")
	s.AddLabel("snap1", "prod")
	s.AddLabel("snap1", "prod")
	labels, _ := s.GetLabels("snap1")
	if len(labels) != 1 {
		t.Fatalf("expected 1 label, got %d", len(labels))
	}
}

func TestRemoveLabel(t *testing.T) {
	s := newLabelStore(t)
	saveTestSnapshot(t, s, "snap1")
	s.AddLabel("snap1", "prod")
	s.AddLabel("snap1", "dev")
	s.RemoveLabel("snap1", "prod")
	labels, _ := s.GetLabels("snap1")
	if len(labels) != 1 || labels[0] != "dev" {
		t.Fatalf("unexpected labels: %v", labels)
	}
}

func TestAddLabelNotFound(t *testing.T) {
	s := newLabelStore(t)
	if err := s.AddLabel("missing", "prod"); err != ErrNotFound {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestListByLabel(t *testing.T) {
	s := newLabelStore(t)
	saveTestSnapshot(t, s, "a")
	saveTestSnapshot(t, s, "b")
	saveTestSnapshot(t, s, "c")
	s.AddLabel("a", "prod")
	s.AddLabel("c", "prod")
	s.AddLabel("b", "dev")
	matched, err := s.ListByLabel("prod")
	if err != nil {
		t.Fatal(err)
	}
	if len(matched) != 2 {
		t.Fatalf("expected 2 matches, got %d", len(matched))
	}
}
