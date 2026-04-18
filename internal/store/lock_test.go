package store

import (
	"testing"
	"time"
)

func TestLockAndUnlock(t *testing.T) {
	s := newTempStore(t)
	unlock, err := s.Lock()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !s.Locked() {
		t.Fatal("expected store to be locked")
	}
	unlock()
	if s.Locked() {
		t.Fatal("expected store to be unlocked after unlock()")
	}
}

func TestLockIsExclusive(t *testing.T) {
	s := newTempStore(t)
	unlock, err := s.Lock()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer unlock()

	done := make(chan error, 1)
	go func() {
		_, err := s.Lock()
		done <- err
	}()

	select {
		case err := <-done:
			if err == nil {
				t.Fatal("expected second lock to fail, but it succeeded")
			}
		case <-time.After(lockTimeout + time.Second):
			t.Fatal("timed out waiting for second lock to fail")
	}
}

func TestLockedReturnsFalseWhenNotLocked(t *testing.T) {
	s := newTempStore(t)
	if s.Locked() {
		t.Fatal("expected store to not be locked initially")
	}
}
