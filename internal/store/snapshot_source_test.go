package store

import (
	"os"
	"path/filepath"
	"testing"
)

func newSourceStore(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	return dir
}

func touchSnapshotForSource(t *testing.T, base, name string) {
	t.Helper()
	f := filepath.Join(base, name+".json")
	if err := os.WriteFile(f, []byte(`{}`), 0600); err != nil {
		t.Fatalf("touch snapshot: %v", err)
	}
}

func TestSetAndGetSource(t *testing.T) {
	base := newSourceStore(t)
	touchSnapshotForSource(t, base, "prod")

	rec := SourceRecord{Hostname: "host1", Directory: "/app", User: "alice"}
	if err := SetSource(base, "prod", rec); err != nil {
		t.Fatalf("SetSource: %v", err)
	}

	got, err := GetSource(base, "prod")
	if err != nil {
		t.Fatalf("GetSource: %v", err)
	}
	if got != rec {
		t.Errorf("got %+v, want %+v", got, rec)
	}
}

func TestGetSourceMissing(t *testing.T) {
	base := newSourceStore(t)
	touchSnapshotForSource(t, base, "prod")

	got, err := GetSource(base, "prod")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != (SourceRecord{}) {
		t.Errorf("expected zero value, got %+v", got)
	}
}

func TestClearSource(t *testing.T) {
	base := newSourceStore(t)
	touchSnapshotForSource(t, base, "prod")

	rec := SourceRecord{Hostname: "host1", Directory: "/app", User: "alice"}
	_ = SetSource(base, "prod", rec)

	if err := ClearSource(base, "prod"); err != nil {
		t.Fatalf("ClearSource: %v", err)
	}

	got, err := GetSource(base, "prod")
	if err != nil {
		t.Fatalf("GetSource after clear: %v", err)
	}
	if got != (SourceRecord{}) {
		t.Errorf("expected zero value after clear, got %+v", got)
	}
}

func TestClearSourceIdempotent(t *testing.T) {
	base := newSourceStore(t)
	touchSnapshotForSource(t, base, "prod")

	if err := ClearSource(base, "prod"); err != nil {
		t.Errorf("ClearSource on missing file should not error: %v", err)
	}
}

func TestSetSourceNotFound(t *testing.T) {
	base := newSourceStore(t)

	rec := SourceRecord{Hostname: "host1"}
	if err := SetSource(base, "ghost", rec); err == nil {
		t.Error("expected error for missing snapshot, got nil")
	}
}
