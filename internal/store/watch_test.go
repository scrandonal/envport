package store

import (
	"testing"
)

func TestSetAndGetWatch(t *testing.T) {
	s := newTempStore(t)
	s.Init("alpha")

	if err := SetWatch(s.Root, "alpha"); err != nil {
		t.Fatalf("SetWatch: %v", err)
	}
	event, err := GetWatch(s.Root, "alpha")
	if err != nil {
		t.Fatalf("GetWatch: %v", err)
	}
	if event == nil {
		t.Fatal("expected event, got nil")
	}
	if event.Name != "alpha" {
		t.Errorf("expected name alpha, got %s", event.Name)
	}
}

func TestGetWatchMissing(t *testing.T) {
	s := newTempStore(t)
	s.Init("alpha")
	event, err := GetWatch(s.Root, "alpha")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if event != nil {
		t.Fatal("expected nil event")
	}
}

func TestClearWatch(t *testing.T) {
	s := newTempStore(t)
	s.Init("alpha")
	SetWatch(s.Root, "alpha")
	if err := ClearWatch(s.Root, "alpha"); err != nil {
		t.Fatalf("ClearWatch: %v", err)
	}
	event, _ := GetWatch(s.Root, "alpha")
	if event != nil {
		t.Fatal("expected nil after clear")
	}
}

func TestClearWatchIdempotent(t *testing.T) {
	s := newTempStore(t)
	s.Init("alpha")
	if err := ClearWatch(s.Root, "alpha"); err != nil {
		t.Fatalf("ClearWatch idempotent: %v", err)
	}
}

func TestSetWatchNotFound(t *testing.T) {
	s := newTempStore(t)
	if err := SetWatch(s.Root, "ghost"); err != ErrNotFound {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestListWatched(t *testing.T) {
	s := newTempStore(t)
	s.Init("alpha")
	s.Init("beta")
	s.Init("gamma")
	SetWatch(s.Root, "alpha")
	SetWatch(s.Root, "gamma")

	events, err := ListWatched(s.Root)
	if err != nil {
		t.Fatalf("ListWatched: %v", err)
	}
	if len(events) != 2 {
		t.Errorf("expected 2 watched, got %d", len(events))
	}
}
