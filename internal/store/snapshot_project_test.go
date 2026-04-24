package store

import (
	"testing"
)

func newProjectStore(t *testing.T) string {
	t.Helper()
	s := newTempStore(t)
	if err := s.Init(); err != nil {
		t.Fatalf("init: %v", err)
	}
	return s.base
}

func TestSetAndGetProject(t *testing.T) {
	base := newProjectStore(t)
	touchSnapshot(t, base, "mysnap")

	info := ProjectInfo{Name: "myproject", URL: "https://github.com/example/myproject"}
	if err := SetProject(base, "mysnap", info); err != nil {
		t.Fatalf("SetProject: %v", err)
	}

	got, err := GetProject(base, "mysnap")
	if err != nil {
		t.Fatalf("GetProject: %v", err)
	}
	if got.Name != info.Name || got.URL != info.URL {
		t.Errorf("got %+v, want %+v", got, info)
	}
}

func TestGetProjectMissing(t *testing.T) {
	base := newProjectStore(t)
	touchSnapshot(t, base, "snap")

	info, err := GetProject(base, "snap")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if info.Name != "" {
		t.Errorf("expected empty, got %+v", info)
	}
}

func TestClearProject(t *testing.T) {
	base := newProjectStore(t)
	touchSnapshot(t, base, "snap")

	_ = SetProject(base, "snap", ProjectInfo{Name: "proj"})
	if err := ClearProject(base, "snap"); err != nil {
		t.Fatalf("ClearProject: %v", err)
	}
	info, _ := GetProject(base, "snap")
	if info.Name != "" {
		t.Errorf("expected cleared, got %+v", info)
	}
}

func TestClearProjectIdempotent(t *testing.T) {
	base := newProjectStore(t)
	touchSnapshot(t, base, "snap")

	if err := ClearProject(base, "snap"); err != nil {
		t.Errorf("expected no error on missing project file, got: %v", err)
	}
}

func TestSetProjectNotFound(t *testing.T) {
	base := newProjectStore(t)
	err := SetProject(base, "ghost", ProjectInfo{Name: "proj"})
	if err != ErrNotFound {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestListByProject(t *testing.T) {
	base := newProjectStore(t)
	touchSnapshot(t, base, "snap1")
	touchSnapshot(t, base, "snap2")
	touchSnapshot(t, base, "snap3")

	_ = SetProject(base, "snap1", ProjectInfo{Name: "alpha"})
	_ = SetProject(base, "snap2", ProjectInfo{Name: "alpha"})
	_ = SetProject(base, "snap3", ProjectInfo{Name: "beta"})

	results, err := ListByProject(base, "alpha")
	if err != nil {
		t.Fatalf("ListByProject: %v", err)
	}
	if len(results) != 2 {
		t.Errorf("expected 2 results, got %d: %v", len(results), results)
	}
}
