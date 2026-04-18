package cmd

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/spf13/cobra"
)

type mockHistoryEntry struct {
	Name      string
	Operation string
	Timestamp time.Time
}

func TestHistoryCmdList(t *testing.T) {
	m := &mockManager{
		history: []mockHistoryEntry{
			{Name: "prod", Operation: "save", Timestamp: time.Now()},
			{Name: "prod", Operation: "load", Timestamp: time.Now()},
		},
	}

	cmd := newHistoryCmd(m)
	buf := &bytes.Buffer{}
	cmd.SetOut(buf)
	cmd.SetErr(buf)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "save") || !strings.Contains(out, "load") {
		t.Errorf("expected operations in output, got: %s", out)
	}
	if !strings.Contains(out, "prod") {
		t.Errorf("expected name in output, got: %s", out)
	}
}

func TestHistoryCmdEmpty(t *testing.T) {
	m := &mockManager{history: []mockHistoryEntry{}}
	cmd := newHistoryCmd(m)
	buf := &bytes.Buffer{}
	cmd.SetOut(buf)
	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "No history") {
		t.Errorf("expected empty message, got: %s", buf.String())
	}
}

func TestHistoryCmdClear(t *testing.T) {
	m := &mockManager{}
	cmd := newHistoryCmd(m)
	buf := &bytes.Buffer{}
	cmd.SetOut(buf)
	_ = executeCmd(cmd, []string{"--clear"})
	if !m.historyCleaned {
		t.Error("expected ClearHistory to be called")
	}
}

func executeCmd(cmd *cobra.Command, args []string) error {
	cmd.SetArgs(args)
	return cmd.Execute()
}
