package cmd

import (
	"bytes"
	"testing"

	"github.com/spf13/cobra"
)

func runTierCmd(t *testing.T, m Manager, args ...string) (string, error) {
	t.Helper()
	root := &cobra.Command{Use: "envport"}
	root.AddCommand(newTierCmd(m))
	buf := &bytes.Buffer{}
	root.SetOut(buf)
	root.SetErr(buf)
	root.SetArgs(args)
	err := root.Execute()
	return buf.String(), err
}

func TestTierCmdSet(t *testing.T) {
	m := &mockManager{}
	out, err := runTierCmd(t, m, "tier", "set", "prod", "premium")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m.tierSet != "premium" || m.tierName != "prod" {
		t.Errorf("expected tier premium for prod, got %q / %q", m.tierSet, m.tierName)
	}
	if out == "" {
		t.Error("expected output")
	}
}

func TestTierCmdGet(t *testing.T) {
	m := &mockManager{tierVal: "standard"}
	out, err := runTierCmd(t, m, "tier", "get", "prod")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != "standard\n" {
		t.Errorf("expected 'standard', got %q", out)
	}
}

func TestTierCmdGetEmpty(t *testing.T) {
	m := &mockManager{tierVal: ""}
	out, err := runTierCmd(t, m, "tier", "get", "prod")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != "no tier set\n" {
		t.Errorf("expected 'no tier set', got %q", out)
	}
}

func TestTierCmdClear(t *testing.T) {
	m := &mockManager{}
	_, err := runTierCmd(t, m, "tier", "clear", "prod")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !m.tierCleared {
		t.Error("expected tier to be cleared")
	}
}

func TestTierCmdListEmpty(t *testing.T) {
	m := &mockManager{tierList: []string{}}
	out, err := runTierCmd(t, m, "tier", "list", "free")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out == "" {
		t.Error("expected output for empty list")
	}
}

func TestTierCmdListNames(t *testing.T) {
	m := &mockManager{tierList: []string{"a", "b"}}
	out, err := runTierCmd(t, m, "tier", "list", "premium")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != "a\nb\n" {
		t.Errorf("unexpected output: %q", out)
	}
}
