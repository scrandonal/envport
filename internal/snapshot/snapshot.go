package snapshot

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// Snapshot represents a named set of environment variables.
type Snapshot struct {
	Name      string            `json:"name"`
	CreatedAt time.Time         `json:"created_at"`
	Env       map[string]string `json:"env"`
}

// New creates a new Snapshot from the provided env map.
func New(name string, env map[string]string) *Snapshot {
	return &Snapshot{
		Name:      name,
		CreatedAt: time.Now().UTC(),
		Env:       env,
	}
}

// Save writes the snapshot as JSON to the given file path.
func (s *Snapshot) Save(path string) error {
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal snapshot: %w", err)
	}
	if err := os.WriteFile(path, data, 0600); err != nil {
		return fmt.Errorf("write snapshot file: %w", err)
	}
	return nil
}

// Load reads a snapshot from a JSON file at the given path.
func Load(path string) (*Snapshot, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read snapshot file: %w", err)
	}
	var s Snapshot
	if err := json.Unmarshal(data, &s); err != nil {
		return nil, fmt.Errorf("unmarshal snapshot: %w", err)
	}
	return &s, nil
}

// FromEnviron parses a slice of "KEY=VALUE" strings into a map.
func FromEnviron(environ []string) map[string]string {
	env := make(map[string]string, len(environ))
	for _, e := range environ {
		for i := 0; i < len(e); i++ {
			if e[i] == '=' {
				env[e[:i]] = e[i+1:]
				break
			}
		}
	}
	return env
}
