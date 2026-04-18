package store

import (
	"testing"
)

func TestHistoryAppendAndRead(t *testing.T) {
	s := newTempStore(t)

	if err := s.AppendHistory("myenv", "save"); err != nil {
		t.Fatalf("AppendHistory: %v", err)
	}
	if err := s.AppendHistory("myenv", "load"); err != nil {
		t.Fatalf("AppendHistory: %v", err)
	}

	entries, err := s.ReadHistory()
	if err != nil {
		t.Fatalf("ReadHistory: %v", err)
	}
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
	if entries[0].Operation != "save" || entries[1].Operation != "load" {
		t.Errorf("unexpected operations: %+v", entries)
	}
}

func TestHistoryEmptyOnMissingFile(t *testing.T) {
	s := newTempStore(t)
	entries, err := s.ReadHistory()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 0 {
		t.Errorf("expected empty history, got %d entries", len(entries))
	}
}

func TestHistoryClear(t *testing.T) {
	s := newTempStore(t)
	_ = s.AppendHistory("env1", "save")

	if err := s.ClearHistory(); err != nil {
		t.Fatalf("ClearHistory: %v", err)
	}
	entries, err := s.ReadHistory()
	if err != nil {
		t.Fatalf("ReadHistory after clear: %v", err)
	}
	if len(entries) != 0 {
		t.Errorf("expected empty history after clear, got %d", len(entries))
	}
}

func TestHistoryClearIdempotent(t *testing.T) {
	s := newTempStore(t)
	if err := s.ClearHistory(); err != nil {
		t.Errorf("ClearHistory on missing file should not error: %v", err)
	}
}
