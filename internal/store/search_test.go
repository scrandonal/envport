package store

import (
	"testing"

	"github.com/nicholasgasior/envport/internal/snapshot"
)

func TestSearchByKey(t *testing.T) {
	m := newManager(t)
	m.Save("dev", &snapshot.Snapshot{Vars: map[string]string{"DEBUG": "true", "PORT": "8080"}})
	m.Save("prod", &snapshot.Snapshot{Vars: map[string]string{"PORT": "443"}})

	results, err := m.Search(SearchOptions{Terms: []string{"DEBUG"}})
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 1 || results[0] != "dev" {
		t.Errorf("expected [dev], got %v", results)
	}
}

func TestSearchByKeyValue(t *testing.T) {
	m := newManager(t)
	m.Save("dev", &snapshot.Snapshot{Vars: map[string]string{"PORT": "8080"}})
	m.Save("prod", &snapshot.Snapshot{Vars: map[string]string{"PORT": "443"}})

	results, err := m.Search(SearchOptions{Terms: []string{"PORT=443"}})
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 1 || results[0] != "prod" {
		t.Errorf("expected [prod], got %v", results)
	}
}

func TestSearchMatchAll(t *testing.T) {
	m := newManager(t)
	m.Save("dev", &snapshot.Snapshot{Vars: map[string]string{"DEBUG": "true", "PORT": "8080"}})
	m.Save("prod", &snapshot.Snapshot{Vars: map[string]string{"PORT": "443"}})

	results, err := m.Search(SearchOptions{Terms: []string{"DEBUG", "PORT"}, MatchAll: true})
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 1 || results[0] != "dev" {
		t.Errorf("expected [dev], got %v", results)
	}
}

func TestSearchNoMatch(t *testing.T) {
	m := newManager(t)
	m.Save("dev", &snapshot.Snapshot{Vars: map[string]string{"PORT": "8080"}})

	results, err := m.Search(SearchOptions{Terms: []string{"MISSING"}})
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 0 {
		t.Errorf("expected no results, got %v", results)
	}
}

func TestSearchOrSemantics(t *testing.T) {
	m := newManager(t)
	m.Save("a", &snapshot.Snapshot{Vars: map[string]string{"FOO": "1"}})
	m.Save("b", &snapshot.Snapshot{Vars: map[string]string{"BAR": "2"}})

	results, err := m.Search(SearchOptions{Terms: []string{"FOO", "BAR"}, MatchAll: false})
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 2 {
		t.Errorf("expected 2 results, got %v", results)
	}
}
