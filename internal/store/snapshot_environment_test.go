package store

import (
	"os"
	"testing"
)

func newEnvironmentStore(t *testing.T) string {
	t.Helper()
	dir, err := os.MkdirTemp("", "envport-env-*")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.RemoveAll(dir) })
	return dir
}

func TestSetAndGetEnvironment(t *testing.T) {
	base := newEnvironmentStore(t)
	touchSnapshot(base, "mysnap")

	rec := EnvironmentRecord{
		Hostname: "myhost",
		User:     "alice",
		OS:       "linux",
		Shell:    "/bin/bash",
	}
	if err := SetEnvironment(base, "mysnap", rec); err != nil {
		t.Fatalf("SetEnvironment: %v", err)
	}

	got, err := GetEnvironment(base, "mysnap")
	if err != nil {
		t.Fatalf("GetEnvironment: %v", err)
	}
	if got.Hostname != rec.Hostname || got.User != rec.User || got.OS != rec.OS || got.Shell != rec.Shell {
		t.Errorf("got %+v, want %+v", got, rec)
	}
}

func TestGetEnvironmentMissing(t *testing.T) {
	base := newEnvironmentStore(t)
	touchSnapshot(base, "snap")

	rec, err := GetEnvironment(base, "snap")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rec.Hostname != "" || rec.User != "" {
		t.Errorf("expected empty record, got %+v", rec)
	}
}

func TestClearEnvironment(t *testing.T) {
	base := newEnvironmentStore(t)
	touchSnapshot(base, "snap")

	_ = SetEnvironment(base, "snap", EnvironmentRecord{Hostname: "h"})
	if err := ClearEnvironment(base, "snap"); err != nil {
		t.Fatalf("ClearEnvironment: %v", err)
	}
	rec, _ := GetEnvironment(base, "snap")
	if rec.Hostname != "" {
		t.Errorf("expected empty after clear, got %+v", rec)
	}
}

func TestClearEnvironmentIdempotent(t *testing.T) {
	base := newEnvironmentStore(t)
	touchSnapshot(base, "snap")

	if err := ClearEnvironment(base, "snap"); err != nil {
		t.Fatalf("ClearEnvironment on missing: %v", err)
	}
}

func TestSetEnvironmentNotFound(t *testing.T) {
	base := newEnvironmentStore(t)
	err := SetEnvironment(base, "ghost", EnvironmentRecord{Hostname: "x"})
	if err == nil {
		t.Fatal("expected error for missing snapshot")
	}
}
