package cmd_test

import (
	"strings"
	"testing"

	"github.com/nicholasgasior/envport/internal/snapshot"
	"github.com/nicholasgasior/envport/internal/store"
)

func TestAccessShowEmpty(t *testing.T) {
	mgr := newManager(t)
	snap := snapshot.New(map[string]string{"A": "1"})
	if err := mgr.Save("mysnap", snap); err != nil {
		t.Fatal(err)
	}

	out := executeCmd(t, mgr, "access", "show", "mysnap")
	if !strings.Contains(out, "load_count:  0") {
		t.Errorf("expected load_count 0, got: %s", out)
	}
	if !strings.Contains(out, "save_count:  0") {
		t.Errorf("expected save_count 0, got: %s", out)
	}
}

func TestAccessShowAfterRecord(t *testing.T) {
	mgr := newManager(t)
	snap := snapshot.New(map[string]string{"X": "y"})
	if err := mgr.Save("proj", snap); err != nil {
		t.Fatal(err)
	}

	root := mgr.Root()
	_ = store.RecordLoad(root, "proj")
	_ = store.RecordLoad(root, "proj")
	_ = store.RecordSave(root, "proj")

	out := executeCmd(t, mgr, "access", "show", "proj")
	if !strings.Contains(out, "load_count:  2") {
		t.Errorf("expected load_count 2, got: %s", out)
	}
	if !strings.Contains(out, "save_count:  1") {
		t.Errorf("expected save_count 1, got: %s", out)
	}
	if !strings.Contains(out, "last_loaded:") {
		t.Errorf("expected last_loaded in output, got: %s", out)
	}
}

func TestAccessClear(t *testing.T) {
	mgr := newManager(t)
	snap := snapshot.New(map[string]string{"K": "v"})
	if err := mgr.Save("env", snap); err != nil {
		t.Fatal(err)
	}

	_ = store.RecordLoad(mgr.Root(), "env")

	out := executeCmd(t, mgr, "access", "clear", "env")
	if !strings.Contains(out, "cleared") {
		t.Errorf("expected cleared message, got: %s", out)
	}

	rec, _ := store.GetAccess(mgr.Root(), "env")
	if rec.LoadCount != 0 {
		t.Errorf("expected 0 after clear, got %d", rec.LoadCount)
	}
}

func TestAccessRequiresArg(t *testing.T) {
	mgr := newManager(t)
	_, err := executeCmdErr(t, mgr, "access", "show")
	if err == nil {
		t.Error("expected error for missing arg")
	}
}
