package store

import (
	"testing"
)

func newRegionStore(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	if err := Init(dir); err != nil {
		t.Fatalf("init: %v", err)
	}
	return dir
}

func touchSnapshotForRegion(t *testing.T, base, name string) {
	t.Helper()
	m := NewManager(base)
	if err := m.Save(name, map[string]string{"K": "V"}); err != nil {
		t.Fatalf("save: %v", err)
	}
}

func TestSetAndGetRegion(t *testing.T) {
	base := newRegionStore(t)
	touchSnapshotForRegion(t, base, "prod")

	if err := SetRegion(base, "prod", "us-east-1"); err != nil {
		t.Fatalf("SetRegion: %v", err)
	}
	r, err := GetRegion(base, "prod")
	if err != nil {
		t.Fatalf("GetRegion: %v", err)
	}
	if r != "us-east-1" {
		t.Errorf("expected us-east-1, got %q", r)
	}
}

func TestGetRegionMissing(t *testing.T) {
	base := newRegionStore(t)
	touchSnapshotForRegion(t, base, "prod")

	r, err := GetRegion(base, "prod")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r != "" {
		t.Errorf("expected empty, got %q", r)
	}
}

func TestClearRegion(t *testing.T) {
	base := newRegionStore(t)
	touchSnapshotForRegion(t, base, "prod")

	_ = SetRegion(base, "prod", "eu-west-1")
	if err := ClearRegion(base, "prod"); err != nil {
		t.Fatalf("ClearRegion: %v", err)
	}
	r, _ := GetRegion(base, "prod")
	if r != "" {
		t.Errorf("expected empty after clear, got %q", r)
	}
}

func TestClearRegionIdempotent(t *testing.T) {
	base := newRegionStore(t)
	touchSnapshotForRegion(t, base, "prod")

	if err := ClearRegion(base, "prod"); err != nil {
		t.Errorf("first clear: %v", err)
	}
	if err := ClearRegion(base, "prod"); err != nil {
		t.Errorf("second clear: %v", err)
	}
}

func TestSetRegionNotFound(t *testing.T) {
	base := newRegionStore(t)
	if err := SetRegion(base, "ghost", "ap-south-1"); err == nil {
		t.Error("expected error for missing snapshot")
	}
}

func TestListByRegion(t *testing.T) {
	base := newRegionStore(t)
	touchSnapshotForRegion(t, base, "prod")
	touchSnapshotForRegion(t, base, "staging")
	touchSnapshotForRegion(t, base, "dev")

	_ = SetRegion(base, "prod", "us-east-1")
	_ = SetRegion(base, "staging", "us-east-1")
	_ = SetRegion(base, "dev", "eu-central-1")

	names, err := ListByRegion(base, "us-east-1")
	if err != nil {
		t.Fatalf("ListByRegion: %v", err)
	}
	if len(names) != 2 {
		t.Errorf("expected 2, got %d: %v", len(names), names)
	}
}
