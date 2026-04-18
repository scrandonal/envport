package store

import (
	"github.com/nicholasgasior/envport/internal/snapshot"
)

// SearchOptions controls how snapshots are matched.
type SearchOptions struct {
	// Terms is a list of "KEY" or "KEY=VALUE" strings.
	Terms []string
	// MatchAll requires every term to match (AND); otherwise OR semantics.
	MatchAll bool
}

// Search returns the names of snapshots that satisfy opts.
func (m *Manager) Search(opts SearchOptions) ([]string, error) {
	names, err := m.List()
	if err != nil {
		return nil, err
	}

	var results []string
	for _, name := range names {
		snap, err := m.Load(name)
		if err != nil {
			continue
		}
		if matchVars(snap.Vars, opts) {
			results = append(results, name)
		}
	}
	return results, nil
}

func matchVars(vars map[string]string, opts SearchOptions) bool {
	for _, term := range opts.Terms {
		matched := evalTerm(vars, term)
		if opts.MatchAll && !matched {
			return false
		}
		if !opts.MatchAll && matched {
			return true
		}
	}
	return opts.MatchAll
}

func evalTerm(vars map[string]string, term string) bool {
	for i, ch := range term {
		if ch == '=' {
			key, val := term[:i], term[i+1:]
			v, ok := vars[key]
			return ok && v == val
		}
	}
	_, ok := vars[term]
	return ok
}

// Ensure Manager satisfies the interface used by snapshot (compile-time check).
var _ interface {
	Search(SearchOptions) ([]string, error)
} = (*Manager)(nil)

// Silence unused import.
var _ = (*snapshot.Snapshot)(nil)
