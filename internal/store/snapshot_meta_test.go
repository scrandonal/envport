package store

import (
	"testing"
)

func TestSetAndGetMeta(t *testing.T) {
	s := newTempStore(t)
	if err := s.Init("proj"); err != nil {
		t.Fatal(err)
	}
	if err := s.SetMeta("proj", "my description"); err != nil {
		t.Fatal(err)
	}
	m, err := s.GetMeta("proj")
	if err != nil {
		t.Fatal(err)
	}
	if m.Description != "my description" {
		t.Errorf("expected 'my description', got %q", m.Description)
	}
	if m.CreatedAt.IsZero() {
		t.Error("expected CreatedAt to be set")
	}
}

func TestGetMetaMissing(t *testing.T) {
	s := newTempStore(t)
	if err := s.Init("proj"); err != nil {
		t.Fatal(err)
	}
	m, err := s.GetMeta("proj")
	if err != nil {
		t.Fatal(err)
	}
	if m.Description != "" {
		t.Error("expected empty meta")
	}
}

func TestClearMeta(t *testing.T) {
	s := newTempStore(t)
	if err := s.Init("proj"); err != nil {
		t.Fatal(err)
	}
	_ = s.SetMeta("proj", "desc")
	if err := s.ClearMeta("proj"); err != nil {
		t.Fatal(err)
	}
	m, _ := s.GetMeta("proj")
	if m.Description != "" {
		t.Error("expected cleared meta")
	}
}

func TestClearMetaIdempotent(t *testing.T) {
	s := newTempStore(t)
	if err := s.Init("proj"); err != nil {
		t.Fatal(err)
	}
	if err := s.ClearMeta("proj"); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestSetMetaNotFound(t *testing.T) {
	s := newTempStore(t)
	err := s.SetMeta("ghost", "desc")
	if err == nil {
		t.Error("expected error for missing snapshot")
	}
}
