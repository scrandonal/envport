package cmd

import (
	"bytes"
	"errors"
	"testing"

	"github.com/nicholasgasior/envport/internal/snapshot"
)

func TestEditCmdNotFound(t *testing.T) {
	m := &mockManager{
		loadFn: func(name string) (*snapshot.Snapshot, error) {
			return nil, errors.New("not found")
		},
	}

	root := NewRootCmd(m)
	root.SetArgs([]string{"edit", "missing"})
	var buf bytes.Buffer
	root.SetErr(&buf)

	err := root.Execute()
	if err == nil {
		t.Fatal("expected error for missing snapshot")
	}
}

func TestEditCmdRequiresArg(t *testing.T) {
	m := &mockManager{}

	root := NewRootCmd(m)
	root.SetArgs([]string{"edit"})
	var buf bytes.Buffer
	root.SetErr(&buf)

	err := root.Execute()
	if err == nil {
		t.Fatal("expected error when no arg provided")
	}
}
