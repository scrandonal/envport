package store

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

type SnapshotStats struct {
	KeyCount  int            `json:"key_count"`
	SizeBytes int64          `json:"size_bytes"`
	KeySizes  map[string]int `json:"key_sizes"`
}

func statsPath(base, name string) string {
	return filepath.Join(base, name+".stats.json")
}

func ComputeStats(s *Store, name string) (*SnapshotStats, error) {
	snap, err := s.manager.Load(name)
	if err != nil {
		return nil, err
	}

	keySizes := make(map[string]int, len(snap.Vars))
	var total int64
	for k, v := range snap.Vars {
		sz := len(k) + len(v)
		keySizes[k] = sz
		total += int64(sz)
	}

	return &SnapshotStats{
		KeyCount:  len(snap.Vars),
		SizeBytes: total,
		KeySizes:  keySizes,
	}, nil
}

func SaveStats(base, name string, stats *SnapshotStats) error {
	data, err := json.MarshalIndent(stats, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(statsPath(base, name), data, 0600)
}

func LoadStats(base, name string) (*SnapshotStats, error) {
	data, err := os.ReadFile(statsPath(base, name))
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, nil
		}
		return nil, err
	}
	var stats SnapshotStats
	if err := json.Unmarshal(data, &stats); err != nil {
		return nil, err
	}
	return &stats, nil
}
