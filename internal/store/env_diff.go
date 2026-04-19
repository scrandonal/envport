package store

import (
	"fmt"
	"sort"

	"github.com/user/envport/internal/snapshot"
)

// EnvDiff represents the difference between a snapshot and the current environment.
type EnvDiff struct {
	Added   map[string]string // in snapshot, not in env
	Removed map[string]string // in env, not in snapshot
	Changed map[string][2]string // key -> [snapshot value, env value]
}

func (d *EnvDiff) HasChanges() bool {
	return len(d.Added)+len(d.Removed)+len(d.Changed) > 0
}

func (d *EnvDiff) Summary() []string {
	var lines []string
	keys := func(m map[string]string) []string {
		var ks []string
		for k := range m {
			ks = append(ks, k)
		}
		sort.Strings(ks)
		return ks
	}
	for _, k := range keys(d.Added) {
		lines = append(lines, fmt.Sprintf("+ %s=%s", k, d.Added[k]))
	}
	for _, k := range keys(d.Removed) {
		lines = append(lines, fmt.Sprintf("- %s=%s", k, d.Removed[k]))
	}
	changedKeys := make([]string, 0, len(d.Changed))
	for k := range d.Changed {
		changedKeys = append(changedKeys, k)
	}
	sort.Strings(changedKeys)
	for _, k := range changedKeys {
		v := d.Changed[k]
		lines = append(lines, fmt.Sprintf("~ %s: %s -> %s", k, v[0], v[1]))
	}
	return lines
}

// DiffWithEnviron compares a named snapshot against the provided environ map.
func (m *Manager) DiffWithEnviron(name string, environ map[string]string) (*EnvDiff, error) {
	snap, err := m.Load(name)
	if err != nil {
		return nil, err
	}
	return diffWithEnviron(snap, environ), nil
}

func diffWithEnviron(snap *snapshot.Snapshot, environ map[string]string) *EnvDiff {
	d := &EnvDiff{
		Added:   make(map[string]string),
		Removed: make(map[string]string),
		Changed: make(map[string][2]string),
	}
	for k, sv := range snap.Vars {
		ev, ok := environ[k]
		if !ok {
			d.Added[k] = sv
		} else if ev != sv {
			d.Changed[k] = [2]string{sv, ev}
		}
	}
	for k, ev := range environ {
		if _, ok := snap.Vars[k]; !ok {
			d.Removed[k] = ev
		}
	}
	return d
}
