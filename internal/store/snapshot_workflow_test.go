package store_test

import (
	"testing"

	"github.com/user/envport/internal/store"
)

func newWorkflowStore(t *testing.T) *store.Store {
	t.Helper()
	s, err := store.New(t.TempDir())
	if err != nil {
		t.Fatalf("store.New: %v", err)
	}
	if err := s.Init(); err != nil {
		t.Fatalf("s.Init: %v", err)
	}
	return s
}

func touchSnapshotForWorkflow(t *testing.T, s *store.Store, name string) {
	t.Helper()
	snap := map[string]string{"KEY": "val"}
	if err := s.Save(name, snap); err != nil {
		t.Fatalf("s.Save(%q): %v", name, err)
	}
}

func TestSaveAndLoadWorkflow(t *testing.T) {
	s := newWorkflowStore(t)
	touchSnapshotForWorkflow(t, s, "dev")

	wf := store.Workflow{
		Name:  "deploy",
		Steps: []string{"build", "test", "push"},
	}
	if err := s.SaveWorkflow("dev", wf); err != nil {
		t.Fatalf("SaveWorkflow: %v", err)
	}

	got, err := s.LoadWorkflow("dev", "deploy")
	if err != nil {
		t.Fatalf("LoadWorkflow: %v", err)
	}
	if got.Name != wf.Name {
		t.Errorf("Name: got %q want %q", got.Name, wf.Name)
	}
	if len(got.Steps) != len(wf.Steps) {
		t.Errorf("Steps len: got %d want %d", len(got.Steps), len(wf.Steps))
	}
}

func TestLoadWorkflowMissing(t *testing.T) {
	s := newWorkflowStore(t)
	touchSnapshotForWorkflow(t, s, "dev")

	_, err := s.LoadWorkflow("dev", "nonexistent")
	if err == nil {
		t.Fatal("expected error for missing workflow")
	}
}

func TestDeleteWorkflow(t *testing.T) {
	s := newWorkflowStore(t)
	touchSnapshotForWorkflow(t, s, "dev")

	wf := store.Workflow{Name: "ci", Steps: []string{"lint", "test"}}
	if err := s.SaveWorkflow("dev", wf); err != nil {
		t.Fatalf("SaveWorkflow: %v", err)
	}
	if err := s.DeleteWorkflow("dev", "ci"); err != nil {
		t.Fatalf("DeleteWorkflow: %v", err)
	}
	_, err := s.LoadWorkflow("dev", "ci")
	if err == nil {
		t.Fatal("expected error after delete")
	}
}

func TestDeleteWorkflowNotFound(t *testing.T) {
	s := newWorkflowStore(t)
	touchSnapshotForWorkflow(t, s, "dev")

	if err := s.DeleteWorkflow("dev", "ghost"); err == nil {
		t.Fatal("expected error deleting nonexistent workflow")
	}
}

func TestListWorkflows(t *testing.T) {
	s := newWorkflowStore(t)
	touchSnapshotForWorkflow(t, s, "dev")

	for _, name := range []string{"alpha", "beta", "gamma"} {
		wf := store.Workflow{Name: name, Steps: []string{"step1"}}
		if err := s.SaveWorkflow("dev", wf); err != nil {
			t.Fatalf("SaveWorkflow(%q): %v", name, err)
		}
	}

	list, err := s.ListWorkflows("dev")
	if err != nil {
		t.Fatalf("ListWorkflows: %v", err)
	}
	if len(list) != 3 {
		t.Errorf("len: got %d want 3", len(list))
	}
}
