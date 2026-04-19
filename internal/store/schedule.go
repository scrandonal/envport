package store

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"time"
)

type Schedule struct {
	Cron      string    `json:"cron"`
	Action    string    `json:"action"`
	CreatedAt time.Time `json:"created_at"`
}

func schedulePath(base, name string) string {
	return filepath.Join(base, name+".schedule.json")
}

func SetSchedule(base, name string, s Schedule) error {
	if _, err := os.Stat(filepath.Join(base, name+".json")); errors.Is(err, os.ErrNotExist) {
		return ErrNotFound
	}
	s.CreatedAt = time.Now()
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(schedulePath(base, name), data, 0600)
}

func GetSchedule(base, name string) (Schedule, error) {
	data, err := os.ReadFile(schedulePath(base, name))
	if errors.Is(err, os.ErrNotExist) {
		return Schedule{}, ErrNotFound
	}
	if err != nil {
		return Schedule{}, err
	}
	var s Schedule
	return s, json.Unmarshal(data, &s)
}

func ClearSchedule(base, name string) error {
	err := os.Remove(schedulePath(base, name))
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	return err
}

func ListScheduled(base string) ([]string, error) {
	entries, err := os.ReadDir(base)
	if errors.Is(err, os.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var names []string
	for _, e := range entries {
		if filepath.Ext(e.Name()) == ".json" {
			base := e.Name()[:len(e.Name())-len(".schedule.json")]
			if len(e.Name()) > len(".schedule.json") && e.Name()[len(e.Name())-len(".schedule.json"):] == ".schedule.json" {
				names = append(names, base)
			}
		}
	}
	return names, nil
}
