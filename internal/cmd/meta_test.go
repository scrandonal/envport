package cmd

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/joshbeard/envport/internal/store"
)

func TestMetaCmdSet(t *testing.T) {
	m := &mockManager{}
	m.setMetaFn = func(name, desc string) error { return nil }
	cmd := newMetaCmd(m)
	cmd.SetArgs([]string{"proj", "my description"})
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}
}

func TestMetaCmdGet(t *testing.T) {
	m := &mockManager{}
	m.getMetaFn = func(name string) (store.Meta, error) {
		return store.Meta{
			Description: "hello world",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}, nil
	}
	var buf bytes.Buffer
	cmd := newMetaCmd(m)
	cmd.SetOut(&buf)
	cmd.SetArgs([]string{"proj"})
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(buf.String(), "hello world") {
		t.Errorf("expected description in output, got: %s", buf.String())
	}
}

func TestMetaCmdGetEmpty(t *testing.T) {
	m := &mockManager{}
	m.getMetaFn = func(name string) (store.Meta, error) {
		return store.Meta{}, nil
	}
	var buf bytes.Buffer
	cmd := newMetaCmd(m)
	cmd.SetOut(&buf)
	cmd.SetArgs([]string{"proj"})
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(buf.String(), "no description") {
		t.Errorf("expected no description message, got: %s", buf.String())
	}
}

func TestMetaCmdClear(t *testing.T) {
	m := &mockManager{}
	m.clearMetaFn = func(name string) error { return nil }
	cmd := newMetaCmd(m)
	cmd.SetArgs([]string{"proj", "--clear"})
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}
}

func TestMetaCmdRequiresArg(t *testing.T) {
	m := &mockManager{}
	cmd := newMetaCmd(m)
	cmd.SetArgs([]string{})
	if err := cmd.Execute(); err == nil {
		t.Error("expected error with no args")
	}
}
