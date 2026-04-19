package cmd

import (
	"strings"
	"testing"

	"envport/internal/store"
)

func TestScheduleCmdSet(t *testing.T) {
	mgr := &mockManager{}
	root := newScheduleCmd(mgr)
	root.SetArgs([]string{"set", "prod", "@daily", "load"})
	out := &strings.Builder{}
	root.SetOut(out)
	if err := root.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out.String(), "prod") {
		t.Errorf("expected output to mention prod, got %q", out.String())
	}
}

func TestScheduleCmdGet(t *testing.T) {
	mgr := &mockManager{
		schedule: map[string]store.Schedule{
			"prod": {Cron: "@daily", Action: "load"},
		},
	}
	root := newScheduleCmd(mgr)
	root.SetArgs([]string{"get", "prod"})
	out := &strings.Builder{}
	root.SetOut(out)
	if err := root.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out.String(), "@daily") {
		t.Errorf("expected @daily in output, got %q", out.String())
	}
}

func TestScheduleCmdClear(t *testing.T) {
	mgr := &mockManager{
		schedule: map[string]store.Schedule{
			"prod": {Cron: "@daily", Action: "load"},
		},
	}
	root := newScheduleCmd(mgr)
	root.SetArgs([]string{"clear", "prod"})
	out := &strings.Builder{}
	root.SetOut(out)
	if err := root.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out.String(), "cleared") {
		t.Errorf("expected cleared in output, got %q", out.String())
	}
}

func TestScheduleCmdRequiresArgs(t *testing.T) {
	mgr := &mockManager{}
	root := newScheduleCmd(mgr)
	root.SetArgs([]string{"set"})
	if err := root.Execute(); err == nil {
		t.Error("expected error for missing args")
	}
}
