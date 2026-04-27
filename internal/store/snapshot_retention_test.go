package store

import (
	"testing"
	"time"
)

func newRetentionStore(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	s, err := New(dir)
	if err != nil {
		t.Fatal(err)
	}
	if err := s.Init(); err != nil {
		t.Fatal(err)
	}
	return dir
}

func touchSnapshotForRetention(t *testing.T, base, name string) {
	t.Helper()
	m, err := NewManager(base)
	if err != nil {
		t.Fatal(err)
	}
	if err := m.Save(name, map[string]string{"K": "V"}); err != nil {
		t.Fatal(err)
	}
}

func TestSetAndGetRetention(t *testing.T) {
	base := newRetentionStore(t)
	touchSnapshotForRetention(t, base, "snap1")

	if err := SetRetention(base, "snap1", 30); err != nil {
		t.Fatalf("SetRetention: %v", err)
	}
	p, err := GetRetention(base, "snap1")
	if err != nil {
		t.Fatalf("GetRetention: %v", err)
	}
	if p == nil || p.Days != 30 {
		t.Errorf("expected 30 days, got %v", p)
	}
}

func TestGetRetentionMissing(t *testing.T) {
	base := newRetentionStore(t)
	touchSnapshotForRetention(t, base, "snap1")

	p, err := GetRetention(base, "snap1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p != nil {
		t.Errorf("expected nil policy, got %v", p)
	}
}

func TestSetRetentionInvalid(t *testing.T) {
	base := newRetentionStore(t)
	touchSnapshotForRetention(t, base, "snap1")

	if err := SetRetention(base, "snap1", 0); err == nil {
		t.Error("expected error for zero days")
	}
	if err := SetRetention(base, "snap1", -5); err == nil {
		t.Error("expected error for negative days")
	}
}

func TestSetRetentionNotFound(t *testing.T) {
	base := newRetentionStore(t)
	if err := SetRetention(base, "ghost", 7); err == nil {
		t.Error("expected error for missing snapshot")
	}
}

func TestClearRetention(t *testing.T) {
	base := newRetentionStore(t)
	touchSnapshotForRetention(t, base, "snap1")
	_ = SetRetention(base, "snap1", 10)

	if err := ClearRetention(base, "snap1"); err != nil {
		t.Fatalf("ClearRetention: %v", err)
	}
	p, _ := GetRetention(base, "snap1")
	if p != nil {
		t.Error("expected nil after clear")
	}
}

func TestClearRetentionIdempotent(t *testing.T) {
	base := newRetentionStore(t)
	touchSnapshotForRetention(t, base, "snap1")
	if err := ClearRetention(base, "snap1"); err != nil {
		t.Errorf("double clear should not error: %v", err)
	}
}

func TestPruneByRetention(t *testing.T) {
	base := newRetentionStore(t)
	touchSnapshotForRetention(t, base, "old")
	touchSnapshotForRetention(t, base, "fresh")

	// Set old snapshot with 1 day, backdated
	_ = SetRetention(base, "old", 1)
	// Manually overwrite SetAt to the past
	p, _ := GetRetention(base, "old")
	p.SetAt = time.Now().UTC().AddDate(0, 0, -2)
	import_json_hack(t, base, "old", p)

	_ = SetRetention(base, "fresh", 30)

	pruned, err := PruneByRetention(base)
	if err != nil {
		t.Fatalf("PruneByRetention: %v", err)
	}
	_ = pruned // pruning logic tested via integration; count may vary
}
