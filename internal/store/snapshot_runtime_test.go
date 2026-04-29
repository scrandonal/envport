package store

import (
	"os"
	"path/filepath"
	"testing"
)

func newRuntimeStore(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	return dir
}

func touchSnapshotForRuntime(t *testing.T, dir, name string) {
	t.Helper()
	f, err := os.Create(filepath.Join(dir, name+".json"))
	if err != nil {
		t.Fatal(err)
	}
	f.Close()
}

func TestSetAndGetRuntime(t *testing.T) {
	dir := newRuntimeStore(t)
	touchSnapshotForRuntime(t, dir, "snap1")
	info := RuntimeInfo{OS: "linux", Arch: "amd64", Host: "myhost", User: "alice", Shell: "bash"}
	if err := SetRuntime(dir, "snap1", info); err != nil {
		t.Fatalf("SetRuntime: %v", err)
	}
	got, err := GetRuntime(dir, "snap1")
	if err != nil {
		t.Fatalf("GetRuntime: %v", err)
	}
	if got.OS != "linux" || got.User != "alice" || got.Shell != "bash" {
		t.Errorf("unexpected runtime: %+v", got)
	}
}

func TestGetRuntimeMissing(t *testing.T) {
	dir := newRuntimeStore(t)
	touchSnapshotForRuntime(t, dir, "snap1")
	got, err := GetRuntime(dir, "snap1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.OS != "" {
		t.Errorf("expected empty RuntimeInfo, got %+v", got)
	}
}

func TestClearRuntime(t *testing.T) {
	dir := newRuntimeStore(t)
	touchSnapshotForRuntime(t, dir, "snap1")
	info := RuntimeInfo{OS: "darwin", Arch: "arm64"}
	_ = SetRuntime(dir, "snap1", info)
	if err := ClearRuntime(dir, "snap1"); err != nil {
		t.Fatalf("ClearRuntime: %v", err)
	}
	got, err := GetRuntime(dir, "snap1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.OS != "" {
		t.Errorf("expected empty after clear, got %+v", got)
	}
}

func TestClearRuntimeIdempotent(t *testing.T) {
	dir := newRuntimeStore(t)
	if err := ClearRuntime(dir, "nonexistent"); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestSetRuntimeNotFound(t *testing.T) {
	dir := newRuntimeStore(t)
	info := RuntimeInfo{OS: "linux"}
	if err := SetRuntime(dir, "ghost", info); err == nil {
		t.Error("expected error for missing snapshot")
	}
}
