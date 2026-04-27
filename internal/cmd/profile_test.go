package cmd

import (
	"bytes"
	"strings"
	"testing"

	"envport/internal/store"
)

func runProfileCmd(t *testing.T, mgr Manager, args ...string) string {
	t.Helper()
	root := newProfileCmd(mgr)
	buf := &bytes.Buffer{}
	root.SetOut(buf)
	root.SetErr(buf)
	root.SetArgs(args)
	_ = root.Execute()
	return buf.String()
}

func TestProfileCmdSet(t *testing.T) {
	mgr := &mockManager{}
	out := runProfileCmd(t, mgr, "set", "dev", "--description", "Dev env", "--author", "alice")
	if !strings.Contains(out, "profile set for") {
		t.Errorf("unexpected output: %q", out)
	}
	if mgr.setProfileCalled != "dev" {
		t.Errorf("expected SetProfile called with 'dev', got %q", mgr.setProfileCalled)
	}
}

func TestProfileCmdGet(t *testing.T) {
	mgr := &mockManager{
		profileResult: store.Profile{Name: "dev", Description: "Dev env", Author: "alice"},
	}
	out := runProfileCmd(t, mgr, "get", "dev")
	if !strings.Contains(out, "alice") {
		t.Errorf("expected author in output, got: %q", out)
	}
}

func TestProfileCmdGetEmpty(t *testing.T) {
	mgr := &mockManager{}
	out := runProfileCmd(t, mgr, "get", "dev")
	if !strings.Contains(out, "no profile set") {
		t.Errorf("expected empty message, got: %q", out)
	}
}

func TestProfileCmdClear(t *testing.T) {
	mgr := &mockManager{}
	out := runProfileCmd(t, mgr, "clear", "dev")
	if !strings.Contains(out, "profile cleared") {
		t.Errorf("unexpected output: %q", out)
	}
}

func TestProfileCmdRequiresArg(t *testing.T) {
	mgr := &mockManager{}
	root := newProfileCmd(mgr)
	root.SetArgs([]string{"set"})
	if err := root.Execute(); err == nil {
		t.Error("expected error for missing argument")
	}
}
