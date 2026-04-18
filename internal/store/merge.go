package store

import (
	"fmt"
)

// MergeStrategy controls how conflicting keys are resolved.
type MergeStrategy int

const (
	// MergeSkip keeps the destination value on conflict.
	MergeSkip MergeStrategy = iota
	// MergeOverwrite replaces the destination value with the source value.
	MergeOverwrite
)

// MergeResult summarises the outcome of a merge operation.
type MergeResult struct {
	Added    []string
	Skipped  []string
	Overwritten []string
}

// Merge combines the vars from src into dst, using the given strategy for
// conflicting keys. Both snapshots must already exist.
func (m *Manager) Merge(src, dst string, strategy MergeStrategy) (MergeResult, error) {
	srcSnap, err := m.Load(src)
	if err != nil {
		return MergeResult{}, fmt.Errorf("merge: load src %q: %w", src, err)
	}
	dstSnap, err := m.Load(dst)
	if err != nil {
		return MergeResult{}, fmt.Errorf("merge: load dst %q: %w", dst, err)
	}

	var result MergeResult
	for k, v := range srcSnap.Vars {
		if _, exists := dstSnap.Vars[k]; exists {
			switch strategy {
			case MergeSkip:
				result.Skipped = append(result.Skipped, k)
				continue
			case MergeOverwrite:
				result.Overwritten = append(result.Overwritten, k)
			}
		} else {
			result.Added = append(result.Added, k)
		}
		dstSnap.Vars[k] = v
	}

	if err := m.Save(dst, dstSnap); err != nil {
		return MergeResult{}, fmt.Errorf("merge: save dst %q: %w", dst, err)
	}
	return result, nil
}
