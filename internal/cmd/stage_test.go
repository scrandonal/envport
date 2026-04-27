package cmd

import (
	"bytes"
	"errors"
	"strings"
	"testing"
)

func runStageCmd(t *testing.T, m Manager, args ...string) (string, error) {
	t.Helper()
	root := newStageCmd(m)
	buf := &bytes.Buffer{}
	root.SetOut(buf)
	root.SetErr(buf)
	root.SetArgs(args)
	err := root.Execute()
	return buf.String(), err
}

func TestStageCmdSet(t *testing.T) {
	m := &mockManager{}
	_, err := runStageCmd(t, m, "set", "mysnap", "staging")
	if err != nil {
		t.Fatal(err)
	}
}

func TestStageCmdGet(t *testing.T) {
	m := &mockManager{stage: "production"}
	out, err := runStageCmd(t, m, "get", "mysnap")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "production") {
		t.Errorf("expected production in output, got %q", out)
	}
}

func TestStageCmdGetEmpty(t *testing.T) {
	m := &mockManager{stage: ""}
	out, err := runStageCmd(t, m, "get", "mysnap")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "(no stage set)") {
		t.Errorf("expected no stage message, got %q", out)
	}
}

func TestStageCmdClear(t *testing.T) {
	m := &mockManager{}
	_, err := runStageCmd(t, m, "clear", "mysnap")
	if err != nil {
		t.Fatal(err)
	}
}

func TestStageCmdList(t *testing.T) {
	m := &mockManager{stageList: []string{"alpha", "beta"}}
	out, err := runStageCmd(t, m, "list", "staging")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "alpha") || !strings.Contains(out, "beta") {
		t.Errorf("expected names in output, got %q", out)
	}
}

func TestStageCmdListEmpty(t *testing.T) {
	m := &mockManager{stageList: []string{}}
	out, err := runStageCmd(t, m, "list", "dev")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "(none)") {
		t.Errorf("expected (none), got %q", out)
	}
}

func TestStageCmdSetError(t *testing.T) {
	m := &mockManager{err: errors.New("invalid stage")}
	_, err := runStageCmd(t, m, "set", "mysnap", "badstage")
	if err == nil {
		t.Error("expected error")
	}
}

func TestStageCmdRequiresTwoArgs(t *testing.T) {
	m := &mockManager{}
	_, err := runStageCmd(t, m, "set", "onlyonearg")
	if err == nil {
		t.Error("expected error for missing second arg")
	}
}
