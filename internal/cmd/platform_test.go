package cmd

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func runPlatformCmd(t *testing.T, m Manager, args ...string) (string, error) {
	t.Helper()
	root := &cobra.Command{Use: "envport"}
	root.AddCommand(newPlatformCmd(m))
	buf := &bytes.Buffer{}
	root.SetOut(buf)
	root.SetErr(buf)
	root.SetArgs(append([]string{"platform"}, args...))
	err := root.Execute()
	return buf.String(), err
}

func TestPlatformCmdSet(t *testing.T) {
	m := &mockManager{}
	out, err := runPlatformCmd(t, m, "set", "mysnap", "linux")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "linux") {
		t.Errorf("expected linux in output, got %q", out)
	}
}

func TestPlatformCmdGet(t *testing.T) {
	m := &mockManager{platform: "darwin"}
	out, err := runPlatformCmd(t, m, "get", "mysnap")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "darwin") {
		t.Errorf("expected darwin in output, got %q", out)
	}
}

func TestPlatformCmdGetEmpty(t *testing.T) {
	m := &mockManager{platform: ""}
	out, err := runPlatformCmd(t, m, "get", "mysnap")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "no platform") {
		t.Errorf("expected no-platform message, got %q", out)
	}
}

func TestPlatformCmdClear(t *testing.T) {
	m := &mockManager{}
	out, err := runPlatformCmd(t, m, "clear", "mysnap")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "cleared") {
		t.Errorf("expected cleared in output, got %q", out)
	}
}

func TestPlatformCmdList(t *testing.T) {
	m := &mockManager{platformList: []string{"snap1", "snap2"}}
	out, err := runPlatformCmd(t, m, "list", "linux")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "snap1") || !strings.Contains(out, "snap2") {
		t.Errorf("expected snap1 and snap2 in output, got %q", out)
	}
}

func TestPlatformCmdListEmpty(t *testing.T) {
	m := &mockManager{platformList: []string{}}
	out, err := runPlatformCmd(t, m, "list", "windows")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "none") {
		t.Errorf("expected none in output, got %q", out)
	}
}

func TestPlatformCmdSetError(t *testing.T) {
	m := &mockManager{err: errors.New("invalid platform")}
	_, err := runPlatformCmd(t, m, "set", "mysnap", "amiga")
	if err == nil {
		t.Error("expected error")
	}
}
