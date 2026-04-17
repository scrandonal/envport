package cmd

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	"github.com/nicholasgasior/envport/internal/snapshot"
)

func TestDiffCmdChanges(t *testing.T) {
	mgr := &mockManager{
		loads: map[string]*snapshot.Snapshot{
			"snap1": {Vars: map[string]string{"A": "1", "B": "old", "C": "only-in-1"}},
			"snap2": {Vars: map[string]string{"A": "1", "B": "new", "D": "only-in-2"}},
		},
	}
	cmd := newDiffCmd(mgr)
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetArgs([]string{"snap1", "snap2"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "~ B: old -> new") {
		t.Errorf("expected changed line, got: %s", out)
	}
	if !strings.Contains(out, "- C=only-in-1") {
		t.Errorf("expected removed line, got: %s", out)
	}
	if !strings.Contains(out, "+ D=only-in-2") {
		t.Errorf("expected added line, got: %s", out)
	}
}

func TestDiffCmdNoDiff(t *testing.T) {
	mgr := &mockManager{
		loads: map[string]*snapshot.Snapshot{
			"snap1": {Vars: map[string]string{"A": "1"}},
			"snap2": {Vars: map[string]string{"A": "1"}},
		},
	}
	cmd := newDiffCmd(mgr)
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetArgs([]string{"snap1", "snap2"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "no differences") {
		t.Errorf("expected no differences message, got: %s", buf.String())
	}
}

func TestDiffCmdNotFound(t *testing.T) {
	mgr := &mockManager{
		loadErr: errors.New("not found"),
	}
	cmd := newDiffCmd(mgr)
	cmd.SetArgs([]string{"missing", "snap2"})
	if err := cmd.Execute(); err == nil {
		t.Fatal("expected error")
	}
}

func TestDiffCmdRequiresTwoArgs(t *testing.T) {
	mgr := &mockManager{}
	cmd := newDiffCmd(mgr)
	cmd.SetArgs([]string{"only-one"})
	if err := cmd.Execute(); err == nil {
		t.Fatal("expected error for missing second arg")
	}
}
