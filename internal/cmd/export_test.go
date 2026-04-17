package cmd

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	"github.com/nicholasgasior/envport/internal/snapshot"
)

func TestExportCmdShell(t *testing.T) {
	m := &mockManager{
		loadFn: func(name string) (*snapshot.Snapshot, error) {
			return &snapshot.Snapshot{Vars: map[string]string{"FOO": "bar", "BAZ": "qux"}}, nil
		},
	}
	cmd := newExportCmd(m)
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetArgs([]string{"mysnap"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "export BAZ=\"") {
		t.Errorf("expected shell export for BAZ, got: %s", out)
	}
	if !strings.Contains(out, "export FOO=\"") {
		t.Errorf("expected shell export for FOO, got: %s", out)
	}
}

func TestExportCmdDotenv(t *testing.T) {
	m := &mockManager{
		loadFn: func(name string) (*snapshot.Snapshot, error) {
			return &snapshot.Snapshot{Vars: map[string]string{"KEY": "value"}}, nil
		},
	}
	cmd := newExportCmd(m)
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetArgs([]string{"mysnap", "--format", "dotenv"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "KEY=value") {
		t.Errorf("expected dotenv line, got: %s", out)
	}
}

func TestExportCmdNotFound(t *testing.T) {
	m := &mockManager{
		loadFn: func(name string) (*snapshot.Snapshot, error) {
			return nil, errors.New("not found")
		},
	}
	cmd := newExportCmd(m)
	cmd.SetArgs([]string{"missing"})
	if err := cmd.Execute(); err == nil {
		t.Fatal("expected error for missing snapshot")
	}
}

func TestExportCmdRequiresArg(t *testing.T) {
	m := &mockManager{}
	cmd := newExportCmd(m)
	cmd.SetArgs([]string{})
	if err := cmd.Execute(); err == nil {
		t.Fatal("expected error when no arg provided")
	}
}
