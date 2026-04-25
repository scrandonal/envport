package store

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

var validTiers = map[string]bool{
	"free":       true,
	"standard":   true,
	"premium":    true,
	"enterprise": true,
}

func tierPath(root, name string) string {
	return filepath.Join(root, name+".tier.json")
}

func SetTier(root, name, tier string) error {
	if !validTiers[tier] {
		return fmt.Errorf("invalid tier %q: must be one of free, standard, premium, enterprise", tier)
	}
	if _, err := os.Stat(filepath.Join(root, name+".json")); errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("snapshot %q not found", name)
	}
	data, err := json.Marshal(tier)
	if err != nil {
		return err
	}
	return os.WriteFile(tierPath(root, name), data, 0600)
}

func GetTier(root, name string) (string, error) {
	data, err := os.ReadFile(tierPath(root, name))
	if errors.Is(err, os.ErrNotExist) {
		return "", nil
	}
	if err != nil {
		return "", err
	}
	var tier string
	if err := json.Unmarshal(data, &tier); err != nil {
		return "", err
	}
	return tier, nil
}

func ClearTier(root, name string) error {
	err := os.Remove(tierPath(root, name))
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	return err
}

func ListByTier(root, tier string) ([]string, error) {
	entries, err := os.ReadDir(root)
	if err != nil {
		return nil, err
	}
	var names []string
	for _, e := range entries {
		if filepath.Ext(e.Name()) != ".json" {
			continue
		}
		base := e.Name()[:len(e.Name())-5]
		if filepath.Ext(base) != "" {
			continue
		}
		t, err := GetTier(root, base)
		if err != nil || t != tier {
			continue
		}
		names = append(names, base)
	}
	return names, nil
}
