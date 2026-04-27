package store

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"time"
)

// RetentionPolicy defines how long a snapshot should be kept.
type RetentionPolicy struct {
	Days    int       `json:"days"`
	SetAt   time.Time `json:"set_at"`
}

func retentionPath(base, name string) string {
	return filepath.Join(base, name+".retention.json")
}

func SetRetention(base, name string, days int) error {
	if !snapshotExists(base, name) {
		return errors.New("snapshot not found: " + name)
	}
	if days <= 0 {
		return errors.New("retention days must be positive")
	}
	p := RetentionPolicy{Days: days, SetAt: time.Now().UTC()}
	data, err := json.Marshal(p)
	if err != nil {
		return err
	}
	return os.WriteFile(retentionPath(base, name), data, 0600)
}

func GetRetention(base, name string) (*RetentionPolicy, error) {
	data, err := os.ReadFile(retentionPath(base, name))
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, nil
		}
		return nil, err
	}
	var p RetentionPolicy
	if err := json.Unmarshal(data, &p); err != nil {
		return nil, err
	}
	return &p, nil
}

func ClearRetention(base, name string) error {
	err := os.Remove(retentionPath(base, name))
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	return err
}

// PruneByRetention removes snapshots whose retention period has elapsed.
func PruneByRetention(base string) ([]string, error) {
	entries, err := os.ReadDir(base)
	if err != nil {
		return nil, err
	}
	var pruned []string
	now := time.Now().UTC()
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		ext := filepath.Ext(e.Name())
		if ext != ".json" {
			continue
		}
		base2 := e.Name()[:len(e.Name())-len(".retention.json")]
		if filepath.Ext(base2) != "" {
			continue
		}
		// Only process plain snapshot files (no extra dots)
		name := e.Name()[:len(e.Name())-len(".json")]
		p, err := GetRetention(base, name)
		if err != nil || p == nil {
			continue
		}
		expiry := p.SetAt.AddDate(0, 0, p.Days)
		if now.After(expiry) {
			_ = os.Remove(filepath.Join(base, e.Name()))
			_ = ClearRetention(base, name)
			pruned = append(pruned, name)
		}
	}
	return pruned, nil
}
