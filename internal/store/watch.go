package store

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

type WatchEvent struct {
	Name      string    `json:"name"`
	ChangedAt time.Time `json:"changed_at"`
}

func watchPath(root, name string) string {
	return filepath.Join(root, name, "watch.json")
}

func SetWatch(root, name string) error {
	if !snapshotExists(root, name) {
		return ErrNotFound
	}
	event := WatchEvent{Name: name, ChangedAt: time.Now()}
	data, err := json.Marshal(event)
	if err != nil {
		return err
	}
	return os.WriteFile(watchPath(root, name), data, 0600)
}

func GetWatch(root, name string) (*WatchEvent, error) {
	data, err := os.ReadFile(watchPath(root, name))
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var event WatchEvent
	if err := json.Unmarshal(data, &event); err != nil {
		return nil, err
	}
	return &event, nil
}

func ClearWatch(root, name string) error {
	err := os.Remove(watchPath(root, name))
	if os.IsNotExist(err) {
		return nil
	}
	return err
}

func ListWatched(root string) ([]WatchEvent, error) {
	entries, err := os.ReadDir(root)
	if err != nil {
		return nil, err
	}
	var results []WatchEvent
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		event, err := GetWatch(root, e.Name())
		if err != nil || event == nil {
			continue
		}
		results = append(results, *event)
	}
	return results, nil
}
