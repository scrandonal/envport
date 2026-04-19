package cmd_test

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	"github.com/nicholasgasior/envport/internal/cmd"
	"github.com/nicholasgasior/envport/internal/snapshot"
)

func TestRollbackCmdSuccess(t *testing.T) {
	m := &mockManager{
		rollbackFn: func(name string, steps int) error { return nil },
		loadFn: func(name string) (*snapshot.Snapshot, error) {
			return snapshot.New(map[string]string{}), nil
		},
	}
	root := cmd.NewRootCmd(m)
	buf := &bytes.Buffer{}
	root.SetOut(buf)
	root.SetArgs([]string{"rollback", "myenv", "2"})
	if err := root.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "myenv") {
		t.Errorf("expected output to mention snapshot name, got: %q", buf.String())
	}
}

func TestRollbackCmdNotFound(t *testing.T) {
	m := &mockManager{
		rollbackFn: func(name string, steps int) error {
			return errors.New("snapshot not found")
		},
	}
	root := cmd.NewRootCmd(m)
	root.SetArgs([]string{"rollback", "ghost"})
	if err := root.Execute(); err == nil {
		t.Error("expected error for missing snapshot")
	}
}

func TestRollbackCmdInvalidSteps(t *testing.T) {
	m := &mockManager{}
	root := cmd.NewRootCmd(m)
	root.SetArgs([]string{"rollback", "myenv", "abc"})
	if err := root.Execute(); err == nil {
		t.Error("expected error for non-integer steps")
	}
}

func TestRollbackCmdRequiresArg(t *testing.T) {
	m := &mockManager{}
	root := cmd.NewRootCmd(m)
	root.SetArgs([]string{"rollback"})
	if err := root.Execute(); err == nil {
		t.Error("expected error when no args provided")
	}
}
