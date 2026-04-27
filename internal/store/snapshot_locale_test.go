package store

import (
	"os"
	"path/filepath"
	"testing"
)

func newLocaleStore(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	return dir
}

func touchSnapshotForLocale(t *testing.T, base, name string) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(base, name+".json"), []byte(`{}`), 0600); err != nil {
		t.Fatal(err)
	}
}

func TestSetAndGetLocale(t *testing.T) {
	base := newLocaleStore(t)
	touchSnapshotForLocale(t, base, "mysnap")
	if err := SetLocale(base, "mysnap", "en_US"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got, err := GetLocale(base, "mysnap")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "en_US" {
		t.Errorf("expected en_US, got %q", got)
	}
}

func TestGetLocaleMissing(t *testing.T) {
	base := newLocaleStore(t)
	got, err := GetLocale(base, "nosnap")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "" {
		t.Errorf("expected empty, got %q", got)
	}
}

func TestSetLocaleInvalid(t *testing.T) {
	base := newLocaleStore(t)
	touchSnapshotForLocale(t, base, "mysnap")
	if err := SetLocale(base, "mysnap", "xx_ZZ"); err == nil {
		t.Error("expected error for invalid locale")
	}
}

func TestSetLocaleNotFound(t *testing.T) {
	base := newLocaleStore(t)
	if err := SetLocale(base, "ghost", "en_US"); err == nil {
		t.Error("expected error for missing snapshot")
	}
}

func TestClearLocale(t *testing.T) {
	base := newLocaleStore(t)
	touchSnapshotForLocale(t, base, "mysnap")
	_ = SetLocale(base, "mysnap", "fr_FR")
	if err := ClearLocale(base, "mysnap"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got, _ := GetLocale(base, "mysnap")
	if got != "" {
		t.Errorf("expected empty after clear, got %q", got)
	}
}

func TestClearLocaleIdempotent(t *testing.T) {
	base := newLocaleStore(t)
	if err := ClearLocale(base, "nosnap"); err != nil {
		t.Errorf("expected no error on double clear, got %v", err)
	}
}

func TestListByLocale(t *testing.T) {
	base := newLocaleStore(t)
	touchSnapshotForLocale(t, base, "a")
	touchSnapshotForLocale(t, base, "b")
	touchSnapshotForLocale(t, base, "c")
	_ = SetLocale(base, "a", "ja_JP")
	_ = SetLocale(base, "b", "ja_JP")
	_ = SetLocale(base, "c", "de_DE")
	names, err := ListByLocale(base, "ja_JP")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(names) != 2 {
		t.Errorf("expected 2 results, got %d", len(names))
	}
}
