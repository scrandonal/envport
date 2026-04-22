package cmd

import (
	"bytes"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func TestEnvironmentCapture(t *testing.T) {
	mgr := &mockManager{}
	cmd := newEnvironmentCmd(mgr)
	cmd.SetOut(&bytes.Buffer{})
	cmd.SetArgs([]string{"capture", "mysnap"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if mgr.setEnvironmentName != "mysnap" {
		t.Errorf("expected setEnvironmentName=mysnap, got %q", mgr.setEnvironmentName)
	}
}

func TestEnvironmentShow(t *testing.T) {
	mgr := &mockManager{
		environmentRecord: mockEnvironmentRecord{Hostname: "testhost", User: "bob", OS: "darwin", Shell: "/bin/zsh"},
	}
	out := &bytes.Buffer{}
	cmd := newEnvironmentCmd(mgr)
	cmd.SetOut(out)
	cmd.SetArgs([]string{"show", "mysnap"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out.String(), "testhost") {
		t.Errorf("expected hostname in output, got: %s", out.String())
	}
	if !strings.Contains(out.String(), "bob") {
		t.Errorf("expected user in output, got: %s", out.String())
	}
}

func TestEnvironmentClear(t *testing.T) {
	mgr := &mockManager{}
	out := &bytes.Buffer{}
	cmd := newEnvironmentCmd(mgr)
	cmd.SetOut(out)
	cmd.SetArgs([]string{"clear", "mysnap"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out.String(), "cleared") {
		t.Errorf("expected 'cleared' in output, got: %s", out.String())
	}
}

func TestEnvironmentCmdRequiresArg(t *testing.T) {
	for _, sub := range []string{"capture", "show", "clear"} {
		t.Run(sub, func(t *testing.T) {
			mgr := &mockManager{}
			cmd := newEnvironmentCmd(mgr)
			cmd.SetOut(&bytes.Buffer{})
			cmd.SetErr(&bytes.Buffer{})
			cmd.SetArgs([]string{sub})
			if err := cmd.Execute(); err == nil {
				t.Errorf("%s: expected error with no args", sub)
			}
		})
	}
}
