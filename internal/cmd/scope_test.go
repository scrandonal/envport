package cmd

import (
	"bytes"
	"errors"
	"testing"

	"github.com/spf13/cobra"
)

func runScopeCmd(t *testing.T, m Manager, args ...string) (string, error) {
	t.Helper()
	root := &cobra.Command{Use: "envport"}
	root.AddCommand(newScopeCmd(m))
	buf := &bytes.Buffer{}
	root.SetOut(buf)
	root.SetErr(buf)
	root.SetArgs(append([]string{"scope"}, args...))
	err := root.Execute()
	return buf.String(), err
}

func TestScopeCmdSet(t *testing.T) {
	m := &mockManager{}
	out, err := runScopeCmd(t, m, "set", "dev", "global")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out == "" {
		t.Error("expected output")
	}
}

func TestScopeCmdGet(t *testing.T) {
	m := &mockManager{scope: "local"}
	out, err := runScopeCmd(t, m, "get", "dev")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != "local\n" {
		t.Errorf("expected 'local', got %q", out)
	}
}

func TestScopeCmdGetEmpty(t *testing.T) {
	m := &mockManager{scope: ""}
	out, err := runScopeCmd(t, m, "get", "dev")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != "(no scope set)\n" {
		t.Errorf("unexpected output: %q", out)
	}
}

func TestScopeCmdClear(t *testing.T) {
	m := &mockManager{}
	_, err := runScopeCmd(t, m, "clear", "dev")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestScopeCmdList(t *testing.T) {
	m := &mockManager{scopeList: []string{"dev", "staging"}}
	out, err := runScopeCmd(t, m, "list", "global")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != "dev\nstaging\n" {
		t.Errorf("unexpected output: %q", out)
	}
}

func TestScopeCmdListEmpty(t *testing.T) {
	m := &mockManager{scopeList: []string{}}
	out, err := runScopeCmd(t, m, "list", "session")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != "(none)\n" {
		t.Errorf("unexpected output: %q", out)
	}
}

func TestScopeCmdNotFound(t *testing.T) {
	m := &mockManager{err: errors.New("snapshot not found")}
	_, err := runScopeCmd(t, m, "set", "ghost", "global")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestScopeCmdRequiresTwoArgs(t *testing.T) {
	m := &mockManager{}
	_, err := runScopeCmd(t, m, "set", "only-one")
	if err == nil {
		t.Fatal("expected error for missing second arg")
	}
}
