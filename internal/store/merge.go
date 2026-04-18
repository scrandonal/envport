package store

import "fmt"

// MergeStrategy controls how key conflicts are resolved.
type MergeStrategy int

const (
	MergeSkip      MergeStrategy = iota // keep destination value
	MergeOverwrite                       // use source value
)

// Merge copies keys from src into dst.
// Existing keys in dst are handled according to strategy.
func (s *Store) Merge(dst, src string, strategy MergeStrategy) error {
	dstSnap, err := s.manager.Load(dst)
	if err != nil {
		return fmt.Errorf("merge: load dst %q: %w", dst, err)
	}

	srcSnap, err := s.manager.Load(src)
	if err != nil {
		return fmt.Errorf("merge: load src %q: %w", src, err)
	}

	if dstSnap.Vars == nil {
		dstSnap.Vars = make(map[string]string)
	}

	for k, v := range srcSnap.Vars {
		if _, exists := dstSnap.Vars[k]; exists && strategy == MergeSkip {
			continue
		}
		dstSnap.Vars[k] = v
	}

	if err := s.manager.Save(dst, dstSnap); err != nil {
		return fmt.Errorf("merge: save dst %q: %w", dst, err)
	}
	return nil
}
