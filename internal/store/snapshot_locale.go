package store

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

var validLocales = map[string]bool{
	"en_US": true, "en_GB": true, "fr_FR": true, "de_DE": true,
	"es_ES": true, "ja_JP": true, "zh_CN": true, "pt_BR": true,
	"it_IT": true, "ko_KR": true, "ru_RU": true, "nl_NL": true,
}

type LocaleRecord struct {
	Locale string `json:"locale"`
}

func localePath(base, name string) string {
	return filepath.Join(base, name+".locale.json")
}

func SetLocale(base, name, locale string) error {
	if _, ok := validLocales[locale]; !ok {
		return fmt.Errorf("invalid locale %q: must be a supported locale tag", locale)
	}
	snap := filepath.Join(base, name+".json")
	if _, err := os.Stat(snap); errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("snapshot %q not found", name)
	}
	rec := LocaleRecord{Locale: locale}
	data, err := json.Marshal(rec)
	if err != nil {
		return err
	}
	return os.WriteFile(localePath(base, name), data, 0600)
}

func GetLocale(base, name string) (string, error) {
	data, err := os.ReadFile(localePath(base, name))
	if errors.Is(err, os.ErrNotExist) {
		return "", nil
	}
	if err != nil {
		return "", err
	}
	var rec LocaleRecord
	if err := json.Unmarshal(data, &rec); err != nil {
		return "", err
	}
	return rec.Locale, nil
}

func ClearLocale(base, name string) error {
	err := os.Remove(localePath(base, name))
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	return err
}

func ListByLocale(base, locale string) ([]string, error) {
	entries, err := os.ReadDir(base)
	if err != nil {
		return nil, err
	}
	var names []string
	for _, e := range entries {
		if filepath.Ext(e.Name()) != ".json" {
			continue
		}
		base2 := e.Name()[:len(e.Name())-5]
		if filepath.Ext(base2) == ".locale" {
			continue
		}
		l, err := GetLocale(base, base2)
		if err == nil && l == locale {
			names = append(names, base2)
		}
	}
	return names, nil
}
