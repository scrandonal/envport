package store

import (
	"testing"
)

func newChecksumStore(t *testing.T) (string, string) {
	t.Helper()
	root := t.TempDir()
	s, err := New(root)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	if err := s.Init(); err != nil {
		t.Fatalf("Init: %v", err)
	}

	m := NewManager(root)
	snap := map[string]string{"FOO": "bar", "BAZ": "qux"}
	if err := m.Save("mysnap", snap); err != nil {
		t.Fatalf("Save: %v", err)
	}
	return root, "mysnap"
}

func TestComputeChecksum(t *testing.T) {
	root, name := newChecksumStore(t)
	sum, err := ComputeChecksum(root, name)
	if err != nil {
		t.Fatalf("ComputeChecksum: %v", err)
	}
	if len(sum) != 64 {
		t.Errorf("expected 64-char hex string, got %q", sum)
	}
}

func TestSaveAndLoadChecksum(t *testing.T) {
	root, name := newChecksumStore(t)
	sum, err := SaveChecksum(root, name)
	if err != nil {
		t.Fatalf("SaveChecksum: %v", err)
	}

	loaded, err := LoadChecksum(root, name)
	if err != nil {
		t.Fatalf("LoadChecksum: %v", err)
	}
	if loaded != sum {
		t.Errorf("expected %q, got %q", sum, loaded)
	}
}

func TestLoadChecksumMissing(t *testing.T) {
	root, _ := newChecksumStore(t)
	sum, err := LoadChecksum(root, "nosnap")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if sum != "" {
		t.Errorf("expected empty string, got %q", sum)
	}
}

func TestVerifyChecksum(t *testing.T) {
	root, name := newChecksumStore(t)
	if _, err := SaveChecksum(root, name); err != nil {
		t.Fatalf("SaveChecksum: %v", err)
	}

	ok, err := VerifyChecksum(root, name)
	if err != nil {
		t.Fatalf("VerifyChecksum: %v", err)
	}
	if !ok {
		t.Error("expected checksum to match")
	}
}

func TestVerifyChecksumNoStored(t *testing.T) {
	root, name := newChecksumStore(t)
	ok, err := VerifyChecksum(root, name)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ok {
		t.Error("expected false when no checksum stored")
	}
}

func TestClearChecksum(t *testing.T) {
	root, name := newChecksumStore(t)
	if _, err := SaveChecksum(root, name); err != nil {
		t.Fatalf("SaveChecksum: %v", err)
	}
	if err := ClearChecksum(root, name); err != nil {
		t.Fatalf("ClearChecksum: %v", err)
	}
	sum, err := LoadChecksum(root, name)
	if err != nil {
		t.Fatalf("LoadChecksum after clear: %v", err)
	}
	if sum != "" {
		t.Errorf("expected empty after clear, got %q", sum)
	}
}

func TestClearChecksumIdempotent(t *testing.T) {
	root, name := newChecksumStore(t)
	if err := ClearChecksum(root, name); err != nil {
		t.Errorf("ClearChecksum on missing file should not error: %v", err)
	}
}
