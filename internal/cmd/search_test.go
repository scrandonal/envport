package cmd

import (
	"bytes"
	"strings"
	"testing"

	"github.com/nicholasgasior/envport/internal/snapshot"
)

func TestSearchCmdKeyMatch(t *testing.T) {
	m := &mockManager{
		list: []string{"dev", "prod"},
		snaps: map[string]*snapshot.Snapshot{
			"dev":  {Vars: map[string]string{"DEBUG": "true", "PORT": "8080"}},
			"prod": {Vars: map[string]string{"PORT": "443"}},
		},
	}
	buf := &bytes.Buffer{}
	cmd := newSearchCmd(m)
	cmd.SetOut(buf)
	cmd.SetArgs([]string{"DEBUG"})
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}
	out := strings.TrimSpace(buf.String())
	if out != "dev" {
		t.Errorf("expected 'dev', got %q", out)
	}
}

func TestSearchCmdKeyValueMatch(t *testing.T) {
	m := &mockManager{
		list: []string{"dev", "prod"},
		snaps: map[string]*snapshot.Snapshot{
			"dev":  {Vars: map[string]string{"PORT": "8080"}},
			"prod": {Vars: map[string]string{"PORT": "443"}},
		},
	}
	buf := &bytes.Buffer{}
	cmd := newSearchCmd(m)
	cmd.SetOut(buf)
	cmd.SetArgs([]string{"PORT=443"})
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}
	out := strings.TrimSpace(buf.String())
	if out != "prod" {
		t.Errorf("expected 'prod', got %q", out)
	}
}

func TestSearchCmdMatchAll(t *testing.T) {
	m := &mockManager{
		list: []string{"dev", "prod"},
		snaps: map[string]*snapshot.Snapshot{
			"dev":  {Vars: map[string]string{"DEBUG": "true", "PORT": "8080"}},
			"prod": {Vars: map[string]string{"PORT": "443"}},
		},
	}
	buf := &bytes.Buffer{}
	cmd := newSearchCmd(m)
	cmd.SetOut(buf)
	cmd.SetArgs([]string{"--all", "DEBUG", "PORT"})
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}
	out := strings.TrimSpace(buf.String())
	if out != "dev" {
		t.Errorf("expected 'dev', got %q", out)
	}
}

func TestSearchCmdNoMatch(t *testing.T) {
	m := &mockManager{
		list: []string{"dev"},
		snaps: map[string]*snapshot.Snapshot{
			"dev": {Vars: map[string]string{"PORT": "8080"}},
		},
	}
	buf := &bytes.Buffer{}
	cmd := newSearchCmd(m)
	cmd.SetOut(buf)
	cmd.SetArgs([]string{"MISSING"})
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}
	if buf.Len() != 0 {
		t.Errorf("expected no output, got %q", buf.String())
	}
}

func TestSearchCmdRequiresArg(t *testing.T) {
	m := &mockManager{}
	cmd := newSearchCmd(m)
	cmd.SetArgs([]string{})
	if err := cmd.Execute(); err == nil {
		t.Error("expected error for missing argument")
	}
}
