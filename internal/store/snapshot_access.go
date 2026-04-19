package store

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

type AccessRecord struct {
	LastLoaded *time.Time `json:"last_loaded,omitempty"`
	LastSaved  *time.Time `json:"last_saved,omitempty"`
	LoadCount  int        `json:"load_count"`
	SaveCount  int        `json:"save_count"`
}

func accessPath(root, name string) string {
	return filepath.Join(root, name+".access.json")
}

func GetAccess(root, name string) (AccessRecord, error) {
	var rec AccessRecord
	data, err := os.ReadFile(accessPath(root, name))
	if os.IsNotExist(err) {
		return rec, nil
	}
	if err != nil {
		return rec, err
	}
	return rec, json.Unmarshal(data, &rec)
}

func RecordLoad(root, name string) error {
	rec, err := GetAccess(root, name)
	if err != nil {
		return err
	}
	now := time.Now().UTC()
	rec.LastLoaded = &now
	rec.LoadCount++
	return saveAccess(root, name, rec)
}

func RecordSave(root, name string) error {
	rec, err := GetAccess(root, name)
	if err != nil {
		return err
	}
	now := time.Now().UTC()
	rec.LastSaved = &now
	rec.SaveCount++
	return saveAccess(root, name, rec)
}

func ClearAccess(root, name string) error {
	err := os.Remove(accessPath(root, name))
	if os.IsNotExist(err) {
		return nil
	}
	return err
}

func saveAccess(root, name string, rec AccessRecord) error {
	data, err := json.MarshalIndent(rec, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(accessPath(root, name), data, 0600)
}
