package cmd

import (
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func TestWatchCmdSet(t *testing.T) {
	m := &mockManager{}
	cmd := newWatchCmd(m)
	cmd.SetOut(&strings.Builder{})
	out := &strings.Builder{}
	cmd.SetOut(out)
	cmd.SetArgs([]string{"set", "alpha"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out.String(), "watching alpha") {
		t.Errorf("unexpected output: %s", out.String())
	}
}

func TestWatchCmdClear(t *testing.T) {
	m := &mockManager{}
	cmd := newWatchCmd(m)
	out := &strings.Builder{}
	cmd.SetOut(out)
	cmd.SetArgs([]string{"clear", "alpha"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out.String(), "cleared watch") {
		t.Errorf("unexpected output: %s", out.String())
	}
}

func TestWatchCmdListEmpty(t *testing.T) {
	m := &mockManager{}
	cmd := newWatchCmd(m)
	out := &strings.Builder{}
	cmd.SetOut(out)
	cmd.SetArgs([]string{"list"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out.String(), "no watched") {
		t.Errorf("unexpected output: %s", out.String())
	}
}

func TestWatchCmdRequiresArg(t *testing.T) {
	m := &mockManager{}
	cmd := newWatchCmd(m)
	cmd.SetOut(&strings.Builder{})
	cmd.SetErr(&strings.Builder{})
	cmd.SetArgs([]string{"set"})
	if err := cmd.Execute(); err == nil {
		t.Fatal("expected error for missing arg")
	}
}

var _ = (*cobra.Command)(nil)
