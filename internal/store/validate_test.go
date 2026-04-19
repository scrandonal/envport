package store

import (
	"errors"
	"testing"

	"github.com/nicholasgasior/envport/internal/snapshot"
)

func newValidateStore(t *testing.T) *Manager {
	t.Helper()
	s := newTempStore(t)
	return NewManager(s)
}

func TestValidateClean(t *testing.T) {
	m := newValidateStore(t)
	snap := snapshot.New(map[string]string{
		"HOME":    "/home/user",
		"_PRIVATE": "yes",
		"GO111MODULE": "on",
	})
	if err := m.Save("clean", snap); err != nil {
		t.Fatal(err)
	}
	if err := m.Validate("clean"); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestValidateInvalidKeys(t *testing.T) {
	m := newValidateStore(t)
	snap := snapshot.New(map[string]string{
		"VALID_KEY": "ok",
		"123BAD":    "nope",
		"also-bad":  "nope",
	})
	if err := m.Save("dirty", snap); err != nil {
		t.Fatal(err)
	}
	err := m.Validate("dirty")
	if err == nil {
		t.Fatal("expected validation error")
	}
	var ve *ValidationError
	if !errors.As(err, &ve) {
		t.Fatalf("expected *ValidationError, got %T", err)
	}
	if len(ve.InvalidKeys) != 2 {
		t.Fatalf("expected 2 invalid keys, got %d: %v", len(ve.InvalidKeys), ve.InvalidKeys)
	}
}

func TestValidateNotFound(t *testing.T) {
	m := newValidateStore(t)
	if err := m.Validate("ghost"); err == nil {
		t.Fatal("expected error for missing snapshot")
	}
}
