package cmd

import (
	"bytes"
	"errors"
	"testing"
)

func TestRenameCmd(t *testing.T) {
	m := &mockManager{
		snapshots: map[string]map[string]string{
			"old": {"FOO": "bar"},
		},
	}
	cmd := newRenameCmd(m)
	buf := &bytes.Buffer{}
	cmd.SetOut(buf)
	cmd.SetArgs([]string{"old", "new"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := m.snapshots["new"]; !ok {
		t.Error("expected snapshot to be renamed to 'new'")
	}
	if _, ok := m.snapshots["old"]; ok {
		t.Error("expected old snapshot to be removed")
	}
	if got := buf.String(); got == "" {
		t.Error("expected output message")
	}
}

func TestRenameCmdNotFound(t *testing.T) {
	m := &mockManager{}
	cmd := newRenameCmd(m)
	cmd.SetArgs([]string{"missing", "new"})
	if err := cmd.Execute(); err == nil {
		t.Fatal("expected error for missing snapshot")
	}
}

func TestRenameCmdRequiresTwoArgs(t *testing.T) {
	m := &mockManager{}
	cmd := newRenameCmd(m)
	cmd.SetArgs([]string{"only-one"})
	if err := cmd.Execute(); err == nil {
		t.Fatal("expected error for wrong arg count")
	}
}
