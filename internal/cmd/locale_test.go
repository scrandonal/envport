package cmd

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func runLocaleCmd(t *testing.T, mgr Manager, args ...string) (string, error) {
	t.Helper()
	root := &cobra.Command{Use: "envport"}
	root.AddCommand(newLocaleCmd(mgr))
	buf := &bytes.Buffer{}
	root.SetOut(buf)
	root.SetErr(buf)
	root.SetArgs(append([]string{"locale"}, args...))
	err := root.Execute()
	return buf.String(), err
}

func TestLocaleCmdSet(t *testing.T) {
	mgr := &mockManager{}
	out, err := runLocaleCmd(t, mgr, "set", "mysnap", "en_US")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "en_US") {
		t.Errorf("expected locale in output, got %q", out)
	}
}

func TestLocaleCmdGet(t *testing.T) {
	mgr := &mockManager{locale: "fr_FR"}
	out, err := runLocaleCmd(t, mgr, "get", "mysnap")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "fr_FR") {
		t.Errorf("expected fr_FR in output, got %q", out)
	}
}

func TestLocaleCmdGetEmpty(t *testing.T) {
	mgr := &mockManager{locale: ""}
	out, err := runLocaleCmd(t, mgr, "get", "mysnap")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "no locale") {
		t.Errorf("expected 'no locale' message, got %q", out)
	}
}

func TestLocaleCmdClear(t *testing.T) {
	mgr := &mockManager{}
	out, err := runLocaleCmd(t, mgr, "clear", "mysnap")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "cleared") {
		t.Errorf("expected 'cleared' in output, got %q", out)
	}
}

func TestLocaleCmdListEmpty(t *testing.T) {
	mgr := &mockManager{localeList: []string{}}
	out, err := runLocaleCmd(t, mgr, "list", "ja_JP")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "no snapshots") {
		t.Errorf("expected 'no snapshots' message, got %q", out)
	}
}

func TestLocaleCmdListNames(t *testing.T) {
	mgr := &mockManager{localeList: []string{"snap1", "snap2"}}
	out, err := runLocaleCmd(t, mgr, "list", "en_US")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "snap1") || !strings.Contains(out, "snap2") {
		t.Errorf("expected snap names in output, got %q", out)
	}
}

func TestLocaleCmdSetError(t *testing.T) {
	mgr := &mockManager{err: errors.New("invalid locale")}
	_, err := runLocaleCmd(t, mgr, "set", "mysnap", "xx_ZZ")
	if err == nil {
		t.Error("expected error")
	}
}

func TestLocaleCmdRequiresTwoArgs(t *testing.T) {
	mgr := &mockManager{}
	_, err := runLocaleCmd(t, mgr, "set", "onlyone")
	if err == nil {
		t.Error("expected error for missing second arg")
	}
}
