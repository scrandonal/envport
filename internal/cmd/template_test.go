package cmd

import (
	"testing"

	"github.com/nicholasgasior/envport/internal/store"
)

func TestTemplateListEmpty(t *testing.T) {
	m := &mockManager{}
	cmd := newTemplateCmd(m)
	out, err := executeCmd(cmd, "list")
	if err != nil {
		t.Fatal(err)
	}
	if !contains(out, "no templates") {
		t.Fatalf("expected empty message, got %q", out)
	}
}

func TestTemplateListNames(t *testing.T) {
	m := &mockManager{templates: []string{"base", "prod"}}
	cmd := newTemplateCmd(m)
	out, err := executeCmd(cmd, "list")
	if err != nil {
		t.Fatal(err)
	}
	if !contains(out, "base") || !contains(out, "prod") {
		t.Fatalf("unexpected output: %q", out)
	}
}

func TestTemplateApply(t *testing.T) {
	m := &mockManager{}
	cmd := newTemplateCmd(m)
	out, err := executeCmd(cmd, "apply", "base", "mysnap")
	if err != nil {
		t.Fatal(err)
	}
	if !contains(out, "mysnap") {
		t.Fatalf("unexpected output: %q", out)
	}
}

func TestTemplateApplyRequiresTwoArgs(t *testing.T) {
	m := &mockManager{}
	cmd := newTemplateCmd(m)
	_, err := executeCmd(cmd, "apply", "only-one")
	if err == nil {
		t.Fatal("expected error for missing arg")
	}
}

func TestTemplateDelete(t *testing.T) {
	m := &mockManager{}
	cmd := newTemplateCmd(m)
	out, err := executeCmd(cmd, "delete", "base")
	if err != nil {
		t.Fatal(err)
	}
	if !contains(out, "deleted") {
		t.Fatalf("unexpected output: %q", out)
	}
}

func TestTemplateDeleteNotFound(t *testing.T) {
	m := &mockManager{failDelete: true}
	cmd := newTemplateCmd(m)
	_, err := executeCmd(cmd, "delete", "ghost")
	if err == nil {
		t.Fatal("expected error")
	}
}

func contains(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(s) > 0 && containsStr(s, sub))
}

func containsStr(s, sub string) bool {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}

var _ = store.Template{}
