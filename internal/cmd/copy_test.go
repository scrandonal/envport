package cmd

import (
	"bytes"
	"errors"
	"testing"

	"github.com/nicholasgasior/envport/internal/snapshot"
)

func TestCopyCmdSuccess(t *testing.T) {
	snap := snapshot.New(map[string]string{"FOO": "bar"})
	mgr := &mockManager{
		loaded: map[string]*snapshot.Snapshot{"prod": snap},
		names:  []string{"prod"},
	}

	root := NewRootCmd(mgr)
	buf := &bytes.Buffer{}
	root.SetOut(buf)
	root.SetArgs([]string{"copy", "prod", "staging"})

	if err := root.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got := buf.String(); got != `Copied "prod" → "staging"
` {
		t.Errorf("unexpected output: %q", got)
	}
}

func TestCopyCmdDestExists(t *testing.T) {
	snap := snapshot.New(map[string]string{"FOO": "bar"})
	mgr := &mockManager{
		loaded: map[string]*snapshot.Snapshot{"prod": snap},
		names:  []string{"prod", "staging"},
	}

	root := NewRootCmd(mgr)
	root.SetArgs([]string{"copy", "prod", "staging"})

	if err := root.Execute(); err == nil {
		t.Fatal("expected error when destination exists")
	}
}

func TestCopyCmdOverwrite(t *testing.T) {
	snap := snapshot.New(map[string]string{"FOO": "bar"})
	mgr := &mockManager{
		loaded: map[string]*snapshot.Snapshot{"prod": snap},
		names:  []string{"prod", "staging"},
	}

	root := NewRootCmd(mgr)
	buf := &bytes.Buffer{}
	root.SetOut(buf)
	root.SetArgs([]string{"copy", "prod", "staging", "--overwrite"})

	if err := root.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCopyCmdSrcNotFound(t *testing.T) {
	mgr := &mockManager{
		loadErr: errors.New("not found"),
	}

	root := NewRootCmd(mgr)
	root.SetArgs([]string{"copy", "missing", "dst"})

	if err := root.Execute(); err == nil {
		t.Fatal("expected error for missing source")
	}
}

func TestCopyCmdRequiresTwoArgs(t *testing.T) {
	mgr := &mockManager{}
	root := NewRootCmd(mgr)
	root.SetArgs([]string{"copy", "only-one"})

	if err := root.Execute(); err == nil {
		t.Fatal("expected error for wrong arg count")
	}
}
