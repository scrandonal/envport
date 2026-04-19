package store

import (
	"testing"

	"github.com/user/envport/internal/snapshot"
)

func TestDiffWithEnvironNoChanges(t *testing.T) {
	snap := &snapshot.Snapshot{Vars: map[string]string{"A": "1", "B": "2"}}
	env := map[string]string{"A": "1", "B": "2"}
	d := diffWithEnviron(snap, env)
	if d.HasChanges() {
		t.Fatal("expected no changes")
	}
}

func TestDiffWithEnvironAdded(t *testing.T) {
	snap := &snapshot.Snapshot{Vars: map[string]string{"A": "1", "B": "2"}}
	env := map[string]string{"A": "1"}
	d := diffWithEnviron(snap, env)
	if _, ok := d.Added["B"]; !ok {
		t.Fatal("expected B in added")
	}
}

func TestDiffWithEnvironRemoved(t *testing.T) {
	snap := &snapshot.Snapshot{Vars: map[string]string{"A": "1"}}
	env := map[string]string{"A": "1", "X": "extra"}
	d := diffWithEnviron(snap, env)
	if _, ok := d.Removed["X"]; !ok {
		t.Fatal("expected X in removed")
	}
}

func TestDiffWithEnvironChanged(t *testing.T) {
	snap := &snapshot.Snapshot{Vars: map[string]string{"A": "old"}}
	env := map[string]string{"A": "new"}
	d := diffWithEnviron(snap, env)
	v, ok := d.Changed["A"]
	if !ok {
		t.Fatal("expected A in changed")
	}
	if v[0] != "old" || v[1] != "new" {
		t.Fatalf("unexpected values: %v", v)
	}
}

func TestDiffWithEnvironSummary(t *testing.T) {
	snap := &snapshot.Snapshot{Vars: map[string]string{"A": "1", "B": "old"}}
	env := map[string]string{"B": "new", "C": "extra"}
	d := diffWithEnviron(snap, env)
	lines := d.Summary()
	if len(lines) != 3 {
		t.Fatalf("expected 3 summary lines, got %d: %v", len(lines), lines)
	}
}

func TestManagerDiffWithEnviron(t *testing.T) {
	s := newTempStore(t)
	m := NewManager(s)
	snap := &snapshot.Snapshot{Vars: map[string]string{"FOO": "bar"}}
	if err := m.Save("test", snap); err != nil {
		t.Fatal(err)
	}
	d, err := m.DiffWithEnviron("test", map[string]string{"FOO": "baz"})
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := d.Changed["FOO"]; !ok {
		t.Fatal("expected FOO in changed")
	}
}

func TestManagerDiffWithEnvironNotFound(t *testing.T) {
	s := newTempStore(t)
	m := NewManager(s)
	_, err := m.DiffWithEnviron("missing", map[string]string{})
	if err == nil {
		t.Fatal("expected error for missing snapshot")
	}
}
