package store

import (
	"os"
	"path/filepath"
	"testing"
)

func newChannelStore(t *testing.T) string {
	t.Helper()
	return t.TempDir()
}

func touchSnapshotForChannel(t *testing.T, base, name string) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(base, name+".json"), []byte(`{}`), 0600); err != nil {
		t.Fatalf("touch: %v", err)
	}
}

func TestSetAndGetChannel(t *testing.T) {
	base := newChannelStore(t)
	touchSnapshotForChannel(t, base, "mysnap")

	if err := SetChannel(base, "mysnap", "beta"); err != nil {
		t.Fatalf("SetChannel: %v", err)
	}
	ch, err := GetChannel(base, "mysnap")
	if err != nil {
		t.Fatalf("GetChannel: %v", err)
	}
	if ch != "beta" {
		t.Errorf("expected beta, got %q", ch)
	}
}

func TestGetChannelMissing(t *testing.T) {
	base := newChannelStore(t)
	ch, err := GetChannel(base, "ghost")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ch != "" {
		t.Errorf("expected empty, got %q", ch)
	}
}

func TestSetChannelInvalid(t *testing.T) {
	base := newChannelStore(t)
	touchSnapshotForChannel(t, base, "mysnap")
	if err := SetChannel(base, "mysnap", "production"); err == nil {
		t.Fatal("expected error for invalid channel")
	}
}

func TestSetChannelNotFound(t *testing.T) {
	base := newChannelStore(t)
	if err := SetChannel(base, "missing", "stable"); err == nil {
		t.Fatal("expected error for missing snapshot")
	}
}

func TestClearChannel(t *testing.T) {
	base := newChannelStore(t)
	touchSnapshotForChannel(t, base, "mysnap")
	_ = SetChannel(base, "mysnap", "alpha")
	if err := ClearChannel(base, "mysnap"); err != nil {
		t.Fatalf("ClearChannel: %v", err)
	}
	ch, _ := GetChannel(base, "mysnap")
	if ch != "" {
		t.Errorf("expected empty after clear, got %q", ch)
	}
}

func TestClearChannelIdempotent(t *testing.T) {
	base := newChannelStore(t)
	if err := ClearChannel(base, "nonexistent"); err != nil {
		t.Fatalf("ClearChannel on missing should be idempotent: %v", err)
	}
}
