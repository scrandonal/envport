package cmd

import (
	"strings"
	"testing"
	"time"

	"github.com/jdoe/envport/internal/store"
)

func TestAuditCmdList(t *testing.T) {
	m := &mockManager{}
	m.auditEntries = []store.AuditEntry{
		{Time: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC), Operation: "save", Name: "dev", Detail: "created"},
	}
	cmd := newAuditCmd(m)
	out, err := executeCmd(cmd)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "save") || !strings.Contains(out, "dev") {
		t.Errorf("expected audit entry in output, got: %s", out)
	}
}

func TestAuditCmdEmpty(t *testing.T) {
	m := &mockManager{}
	cmd := newAuditCmd(m)
	out, err := executeCmd(cmd)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "no audit entries") {
		t.Errorf("expected empty message, got: %s", out)
	}
}

func TestAuditCmdClear(t *testing.T) {
	m := &mockManager{}
	cmd := newAuditCmd(m)
	cmd.SetArgs([]string{"--clear"})
	out, err := executeCmd(cmd)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "cleared") {
		t.Errorf("expected cleared message, got: %s", out)
	}
	if !m.auditCleared {
		t.Error("expected ClearAuditLog to be called")
	}
}
