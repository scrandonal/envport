package cmd

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func runStatusCmd(t *testing.T, m Manager, args ...string) (string, error) {
	t.Helper()
	root := &cobra.Command{Use: "envport"}
	root.AddCommand(newStatusCmd(m))
	buf := &bytes.Buffer{}
	root.SetOut(buf)
	root.SetErr(buf)
	root.SetArgs(append([]string{"status"}, args...))
	err := root.Execute()
	return buf.String(), err
}

func TestStatusCmdSet(t *testing.T) {
	m := &mockManager{}
	m.setStatusFn = func(name, status string) error { return nil }
	out, err := runStatusCmd(t, m, "set", "prod", "active")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "active") {
		t.Errorf("expected active in output, got %q", out)
	}
}

func TestStatusCmdGet(t *testing.T) {
	m := &mockManager{}
	m.getStatusFn = func(name string) (string, error) { return "draft", nil }
	out, err := runStatusCmd(t, m, "get", "prod")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "draft") {
		t.Errorf("expected draft in output, got %q", out)
	}
}

func TestStatusCmdGetEmpty(t *testing.T) {
	m := &mockManager{}
	m.getStatusFn = func(name string) (string, error) { return "", nil }
	out, err := runStatusCmd(t, m, "get", "prod")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "no status") {
		t.Errorf("expected 'no status' in output, got %q", out)
	}
}

func TestStatusCmdClear(t *testing.T) {
	m := &mockManager{}
	m.clearStatusFn = func(name string) error { return nil }
	out, err := runStatusCmd(t, m, "clear", "prod")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "cleared") {
		t.Errorf("expected cleared in output, got %q", out)
	}
}

func TestStatusCmdList(t *testing.T) {
	m := &mockManager{}
	m.listByStatusFn = func(status string) ([]string, error) { return []string{"prod", "staging"}, nil }
	out, err := runStatusCmd(t, m, "list", "active")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "prod") || !strings.Contains(out, "staging") {
		t.Errorf("unexpected output: %q", out)
	}
}

func TestStatusCmdSetNotFound(t *testing.T) {
	m := &mockManager{}
	m.setStatusFn = func(name, status string) error { return errors.New("not found") }
	_, err := runStatusCmd(t, m, "set", "ghost", "active")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestStatusCmdRequiresTwoArgs(t *testing.T) {
	m := &mockManager{}
	_, err := runStatusCmd(t, m, "set", "onlyone")
	if err == nil {
		t.Fatal("expected error for missing second arg")
	}
}
