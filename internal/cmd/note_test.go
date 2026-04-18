package cmd

import (
	"bytes"
	"strings"
	"testing"

	"envport/internal/store"
)

func TestNoteCmdSet(t *testing.T) {
	m := &mockManager{}
	m.setNoteFunc = func(name, text string) error { return nil }
	cmd := newNoteCmd(m)
	cmd.SetArgs([]string{"proj", "hello world"})
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}
}

func TestNoteCmdGet(t *testing.T) {
	m := &mockManager{}
	m.getNoteFunc = func(name string) (string, error) { return "my note", nil }
	buf := &bytes.Buffer{}
	cmd := newNoteCmd(m)
	cmd.SetOut(buf)
	cmd.SetArgs([]string{"proj"})
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(buf.String(), "my note") {
		t.Errorf("expected note in output, got %q", buf.String())
	}
}

func TestNoteCmdGetEmpty(t *testing.T) {
	m := &mockManager{}
	m.getNoteFunc = func(name string) (string, error) { return "", store.ErrNoteNotFound }
	buf := &bytes.Buffer{}
	cmd := newNoteCmd(m)
	cmd.SetOut(buf)
	cmd.SetArgs([]string{"proj"})
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(buf.String(), "(no note)") {
		t.Errorf("expected '(no note)', got %q", buf.String())
	}
}

func TestNoteCmdClear(t *testing.T) {
	m := &mockManager{}
	m.clearNoteFunc = func(name string) error { return nil }
	cmd := newNoteCmd(m)
	cmd.SetArgs([]string{"proj", "--clear"})
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}
}
