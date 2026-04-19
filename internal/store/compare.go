package store

import "github.com/nicholasgasior/envport/internal/snapshot"

// Diff represents the difference between two snapshots.
type Diff struct {
	Added   map[string]string
	Removed map[string]string
	Changed map[string]DiffChange
}

// DiffChange holds the before and after values for a changed key.
type DiffChange struct {
	Old string
	New string
}

// Compare returns the diff between two named snapshots.
func (m *Manager) Compare(srcName, dstName string) (*Diff, error) {
	src, err := m.Load(srcName)
	if err != nil {
		return nil, err
	}
	dst, err := m.Load(dstName)
	if err != nil {
		return nil, err
	}
	return diffSnapshots(src, dst), nil
}

func diffSnapshots(src, dst *snapshot.Snapshot) *Diff {
	d := &Diff{
		Added:   make(map[string]string),
		Removed: make(map[string]string),
		Changed: make(map[string]DiffChange),
	}
	for k, v := range dst.Vars {
		if old, ok := src.Vars[k]; !ok {
			d.Added[k] = v
		} else if old != v {
			d.Changed[k] = DiffChange{Old: old, New: v}
		}
	}
	for k, v := range src.Vars {
		if _, ok := dst.Vars[k]; !ok {
			d.Removed[k] = v
		}
	}
	return d
}
