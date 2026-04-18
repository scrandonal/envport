package cmd

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func TestLockCmdNotLocked(t *testing.T) {
	m := &mockManager{isLocked: false}
	root := &cobra.Command{Use: "root"}
	root.AddCommand(newLockCmd(m))
	var buf bytes.Buffer
	root.SetOut(&buf)
	root.SetArgs([]string{"lock"})
	if err := root.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "not locked") {
		t.Errorf("expected 'not locked', got %q", buf.String())
	}
}

func TestLockCmdLocked(t *testing.T) {
	m := &mockManager{isLocked: true}
	root := &cobra.Command{Use: "root"}
	root.AddCommand(newLockCmd(m))
	var buf bytes.Buffer
	root.SetOut(&buf)
	root.SetArgs([]string{"lock"})
	if err := root.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "store is locked") {
		t.Errorf("expected 'store is locked', got %q", buf.String())
	}
}

func TestLockCmdForceUnlock(t *testing.T) {
	m := &mockManager{isLocked: true}
	root := &cobra.Command{Use: "root"}
	root.AddCommand(newLockCmd(m))
	var buf bytes.Buffer
	root.SetOut(&buf)
	root.SetArgs([]string{"lock", "--unlock"})
	if err := root.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "lock released") {
		t.Errorf("expected 'lock released', got %q", buf.String())
	}
}

func TestLockCmdForceUnlockError(t *testing.T) {
	m := &mockManager{forceUnlockErr: errors.New("no lock file")}
	root := &cobra.Command{Use: "root"}
	root.AddCommand(newLockCmd(m))
	root.SetArgs([]string{"lock", "--unlock"})
	if err := root.Execute(); err == nil {
		t.Fatal("expected error, got nil")
	}
}
