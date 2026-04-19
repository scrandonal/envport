package store

import (
	"testing"
)

func newScheduleStore(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	return dir
}

func TestSetAndGetSchedule(t *testing.T) {
	dir := newScheduleStore(t)
	createSnapshot(t, dir, "prod")

	s := Schedule{Cron: "0 * * * *", Action: "load"}
	if err := SetSchedule(dir, "prod", s); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got, err := GetSchedule(dir, "prod")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Cron != s.Cron || got.Action != s.Action {
		t.Errorf("got %+v, want %+v", got, s)
	}
}

func TestGetScheduleMissing(t *testing.T) {
	dir := newScheduleStore(t)
	createSnapshot(t, dir, "prod")

	_, err := GetSchedule(dir, "prod")
	if err != ErrNotFound {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestClearSchedule(t *testing.T) {
	dir := newScheduleStore(t)
	createSnapshot(t, dir, "prod")

	_ = SetSchedule(dir, "prod", Schedule{Cron: "@daily", Action: "load"})
	if err := ClearSchedule(dir, "prod"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_, err := GetSchedule(dir, "prod")
	if err != ErrNotFound {
		t.Errorf("expected ErrNotFound after clear, got %v", err)
	}
}

func TestClearScheduleIdempotent(t *testing.T) {
	dir := newScheduleStore(t)
	if err := ClearSchedule(dir, "ghost"); err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestSetScheduleNotFound(t *testing.T) {
	dir := newScheduleStore(t)
	err := SetSchedule(dir, "missing", Schedule{Cron: "@hourly"})
	if err != ErrNotFound {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}
