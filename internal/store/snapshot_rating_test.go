package store

import (
	"testing"
)

func newRatingStore(t *testing.T) string {
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

func TestSetAndGetRating(t *testing.T) {
	dir := newRatingStore(t)
	touchSnapshot(t, dir, "prod")

	if err := SetRating(dir, "prod", 4, "solid"); err != nil {
		t.Fatalf("SetRating: %v", err)
	}
	r, err := GetRating(dir, "prod")
	if err != nil {
		t.Fatalf("GetRating: %v", err)
	}
	if r == nil {
		t.Fatal("expected rating, got nil")
	}
	if r.Value != 4 {
		t.Errorf("expected value 4, got %d", r.Value)
	}
	if r.Comment != "solid" {
		t.Errorf("expected comment 'solid', got %q", r.Comment)
	}
}

func TestGetRatingMissing(t *testing.T) {
	dir := newRatingStore(t)
	touchSnapshot(t, dir, "prod")

	r, err := GetRating(dir, "prod")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r != nil {
		t.Errorf("expected nil rating, got %+v", r)
	}
}

func TestSetRatingInvalid(t *testing.T) {
	dir := newRatingStore(t)
	touchSnapshot(t, dir, "prod")

	if err := SetRating(dir, "prod", 0, ""); err == nil {
		t.Error("expected error for rating 0")
	}
	if err := SetRating(dir, "prod", 6, ""); err == nil {
		t.Error("expected error for rating 6")
	}
}

func TestSetRatingNotFound(t *testing.T) {
	dir := newRatingStore(t)
	if err := SetRating(dir, "ghost", 3, ""); err == nil {
		t.Error("expected error for missing snapshot")
	}
}

func TestClearRating(t *testing.T) {
	dir := newRatingStore(t)
	touchSnapshot(t, dir, "prod")

	_ = SetRating(dir, "prod", 5, "best")
	if err := ClearRating(dir, "prod"); err != nil {
		t.Fatalf("ClearRating: %v", err)
	}
	r, _ := GetRating(dir, "prod")
	if r != nil {
		t.Error("expected nil after clear")
	}
}

func TestClearRatingIdempotent(t *testing.T) {
	dir := newRatingStore(t)
	touchSnapshot(t, dir, "prod")
	if err := ClearRating(dir, "prod"); err != nil {
		t.Errorf("ClearRating on missing should not error: %v", err)
	}
}
