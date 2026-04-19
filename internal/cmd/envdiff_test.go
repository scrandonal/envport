package cmd

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/envport/internal/snapshot"
	"github.com/user/envport/internal/store"
)

func TestEnvDiffNoDifferences(t *testing.T) {
	m := &mockManager{
		diffWithEnvironFn: func(name string, env map[string]string) (*store.EnvDiff, error) {
			return &store.EnvDiff{
				Added:   map[string]string{},
				Removed: map[string]string{},
				Changed: map[string][2]string{},
			}, nil
		},
	}
	cmd := newEnvDiffCmd(m)
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetArgs([]string{"mysnap"})
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(buf.String(), "no differences") {
		t.Fatalf("unexpected output: %q", buf.String())
	}
}

func TestEnvDiffShowsChanges(t *testing.T) {
	m := &mockManager{
		diffWithEnvironFn: func(name string, env map[string]string) (*store.EnvDiff, error) {
			return &store.EnvDiff{
				Added:   map[string]string{"MISSING": "val"},
				Removed: map[string]string{},
				Changed: map[string][2]string{},
			}, nil
		},
	}
	cmd := newEnvDiffCmd(m)
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetArgs([]string{"mysnap"})
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(buf.String(), "+ MISSING") {
		t.Fatalf("unexpected output: %q", buf.String())
	}
}

func TestEnvDiffNotFound(t *testing.T) {
	m := &mockManager{
		diffWithEnvironFn: func(name string, env map[string]string) (*store.EnvDiff, error) {
			return nil, snapshot.ErrNotFound
		},
	}
	cmd := newEnvDiffCmd(m)
	cmd.SetArgs([]string{"missing"})
	if err := cmd.Execute(); err == nil {
		t.Fatal("expected error")
	}
}

func TestEnvDiffRequiresArg(t *testing.T) {
	m := &mockManager{}
	cmd := newEnvDiffCmd(m)
	if err := cmd.Execute(); err == nil {
		t.Fatal("expected error without arg")
	}
}
