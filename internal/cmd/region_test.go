package cmd

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func runRegionCmd(m Manager, args ...string) (string, error) {
	root := &cobra.Command{Use: "envport"}
	root.AddCommand(newRegionCmd(m))
	buf := &bytes.Buffer{}
	root.SetOut(buf)
	root.SetErr(buf)
	root.SetArgs(append([]string{"region"}, args...))
	err := root.Execute()
	return buf.String(), err
}

func TestRegionCmdSet(t *testing.T) {
	m := &mockManager{}
	out, err := runRegionCmd(m, "set", "prod", "us-east-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "us-east-1") {
		t.Errorf("expected region in output, got: %q", out)
	}
}

func TestRegionCmdGet(t *testing.T) {
	m := &mockManager{region: "eu-west-1"}
	out, err := runRegionCmd(m, "get", "prod")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "eu-west-1") {
		t.Errorf("expected region in output, got: %q", out)
	}
}

func TestRegionCmdGetEmpty(t *testing.T) {
	m := &mockManager{region: ""}
	out, err := runRegionCmd(m, "get", "prod")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "no region") {
		t.Errorf("expected 'no region' message, got: %q", out)
	}
}

func TestRegionCmdClear(t *testing.T) {
	m := &mockManager{}
	out, err := runRegionCmd(m, "clear", "prod")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "cleared") {
		t.Errorf("expected 'cleared' in output, got: %q", out)
	}
}

func TestRegionCmdList(t *testing.T) {
	m := &mockManager{regionList: []string{"prod", "staging"}}
	out, err := runRegionCmd(m, "list", "us-east-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "prod") || !strings.Contains(out, "staging") {
		t.Errorf("expected snapshot names in output, got: %q", out)
	}
}

func TestRegionCmdListEmpty(t *testing.T) {
	m := &mockManager{regionList: []string{}}
	out, err := runRegionCmd(m, "list", "ap-south-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "no snapshots") {
		t.Errorf("expected 'no snapshots' message, got: %q", out)
	}
}

func TestRegionCmdNotFound(t *testing.T) {
	m := &mockManager{err: errors.New("snapshot \"ghost\" not found")}
	_, err := runRegionCmd(m, "set", "ghost", "us-west-2")
	if err == nil {
		t.Error("expected error for missing snapshot")
	}
}

func TestRegionCmdRequiresTwoArgs(t *testing.T) {
	m := &mockManager{}
	_, err := runRegionCmd(m, "set", "only-one")
	if err == nil {
		t.Error("expected error for missing second arg")
	}
}
