package store

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

var validChannels = map[string]bool{
	"stable":  true,
	"beta":    true,
	"alpha":   true,
	"nightly": true,
	"canary":  true,
}

func channelPath(base, name string) string {
	return filepath.Join(base, name+".channel.json")
}

// SetChannel assigns a release channel to a snapshot.
func SetChannel(base, name, channel string) error {
	if !validChannels[channel] {
		return fmt.Errorf("invalid channel %q: must be one of stable, beta, alpha, nightly, canary", channel)
	}
	if _, err := os.Stat(filepath.Join(base, name+".json")); errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("snapshot %q not found", name)
	}
	data, _ := json.Marshal(channel)
	return os.WriteFile(channelPath(base, name), data, 0600)
}

// GetChannel returns the channel for a snapshot, or empty string if unset.
func GetChannel(base, name string) (string, error) {
	data, err := os.ReadFile(channelPath(base, name))
	if errors.Is(err, os.ErrNotExist) {
		return "", nil
	}
	if err != nil {
		return "", err
	}
	var ch string
	if err := json.Unmarshal(data, &ch); err != nil {
		return "", err
	}
	return ch, nil
}

// ClearChannel removes the channel assignment from a snapshot.
func ClearChannel(base, name string) error {
	err := os.Remove(channelPath(base, name))
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	return err
}

// ListByChannel returns all snapshot names assigned to the given channel.
func ListByChannel(base, channel string) ([]string, error) {
	entries, err := os.ReadDir(base)
	if err != nil {
		return nil, err
	}
	var names []string
	for _, e := range entries {
		if filepath.Ext(e.Name()) != ".json" {
			continue
		}
		name := e.Name()[:len(e.Name())-len(".channel.json")]
		if len(e.Name()) <= len(".channel.json") || e.Name()[len(e.Name())-len(".channel.json"):] != ".channel.json" {
			continue
		}
		name = e.Name()[:len(e.Name())-len(".channel.json")]
		ch, err := GetChannel(base, name)
		if err != nil || ch != channel {
			continue
		}
		names = append(names, name)
	}
	return names, nil
}
