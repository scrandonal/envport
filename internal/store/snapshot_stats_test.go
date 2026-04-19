package store_test

import (
	"testing"

	"github.com/nicholasgasior/envport/internal/snapshot"
	"github.com/nicholasgasior/envport/internal/store"
)

func newStatsStore(t *testing.T) *store.Store {
	t.Helper()
	s := newTempStore(t)
	return s
}

func TestComputeStats(t *testing.T) {
	s := newStatsStore(t)
	snap := snapshot.New(map[string]string{"FOO": "bar", "BAZ": "qux"})
	if err := s.Manager().Save("mysnap", snap); err != nil {
		t.Fatalf("save: %v", err)
	}

	stats, err := store.ComputeStats(s, "mysnap")
	if err != nil {
		t.Fatalf("compute: %v", err)
	}
	if stats.KeyCount != 2 {
		t.Errorf("expected 2 keys, got %d", stats.KeyCount)
	}
	if stats.SizeBytes == 0 {
		t.Error("expected non-zero size")
	}
	if _, ok := stats.KeySizes["FOO"]; !ok {
		t.Error("expected FOO in key sizes")
	}
}

func TestSaveAndLoadStats(t *testing.T) {
	dir := t.TempDir()
	stats := &store.SnapshotStats{
		KeyCount:  3,
		SizeBytes: 42,
		KeySizes:  map[string]int{"A": 10, "B": 20, "C": 12},
	}
	if err := store.SaveStats(dir, "snap1", stats); err != nil {
		t.Fatalf("save: %v", err)
	}
	loaded, err := store.LoadStats(dir, "snap1")
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if loaded.KeyCount != 3 {
		t.Errorf("expected 3, got %d", loaded.KeyCount)
	}
	if loaded.SizeBytes != 42 {
		t.Errorf("expected 42, got %d", loaded.SizeBytes)
	}
}

func TestLoadStatsMissing(t *testing.T) {
	dir := t.TempDir()
	result, err := store.LoadStats(dir, "nonexistent")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != nil {
		t.Error("expected nil for missing stats")
	}
}
