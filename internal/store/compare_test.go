package store

import (
	"testing"

	"github.com/nicholasgasior/envport/internal/snapshot"
)

func TestDiffSnapshots(t *testing.T) {
	src := &snapshot.Snapshot{Vars: map[string]string{
		"A": "1",
		"B": "2",
		"C": "3",
	}}
	dst := &snapshot.Snapshot{Vars: map[string]string{
		"A": "1",
		"B": "changed",
		"D": "4",
	}}
	d := diffSnapshots(src, dst)

	if _, ok := d.Added["D"]; !ok {
		t.Error("expected D to be added")
	}
	if _, ok := d.Removed["C"]; !ok {
		t.Error("expected C to be removed")
	}
	if ch, ok := d.Changed["B"]; !ok || ch.Old != "2" || ch.New != "changed" {
		t.Errorf("unexpected change for B: %+v", ch)
	}
	if _, ok := d.Changed["A"]; ok {
		t.Error("A should not appear in changed")
	}
}

func TestManagerCompare(t *testing.T) {
	m := newManager(t)

	src := &snapshot.Snapshot{Vars: map[string]string{"X": "1", "Y": "2"}}
	dst := &snapshot.Snapshot{Vars: map[string]string{"X": "99", "Z": "3"}}

	if err := m.Save("snap-a", src); err != nil {
		t.Fatal(err)
	}
	if err := m.Save("snap-b", dst); err != nil {
		t.Fatal(err)
	}

	d, err := m.Compare("snap-a", "snap-b")
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := d.Added["Z"]; !ok {
		t.Error("expected Z added")
	}
	if _, ok := d.Removed["Y"]; !ok {
		t.Error("expected Y removed")
	}
	if ch, ok := d.Changed["X"]; !ok || ch.Old != "1" || ch.New != "99" {
		t.Errorf("unexpected X change: %+v", ch)
	}
}

func TestManagerCompareNotFound(t *testing.T) {
	m := newManager(t)
	_, err := m.Compare("missing", "also-missing")
	if err == nil {
		t.Error("expected error for missing snapshots")
	}
}
