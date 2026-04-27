package store

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// WorkflowStep represents a single step in a snapshot workflow.
type WorkflowStep struct {
	Name    string            `json:"name"`
	Action  string            `json:"action"` // e.g. "load", "merge", "copy", "validate"
	Target  string            `json:"target,omitempty"`
	Options map[string]string `json:"options,omitempty"`
}

// Workflow defines an ordered sequence of steps to execute against snapshots.
type Workflow struct {
	Name        string         `json:"name"`
	Description string         `json:"description,omitempty"`
	Steps       []WorkflowStep `json:"steps"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}

// workflowPath returns the path to the workflow definition file.
func workflowPath(root, name string) string {
	return filepath.Join(root, name+".workflow.json")
}

// SaveWorkflow persists a workflow definition to the store directory.
func SaveWorkflow(root, name string, wf Workflow) error {
	if name == "" {
		return errors.New("workflow name must not be empty")
	}
	if len(wf.Steps) == 0 {
		return errors.New("workflow must contain at least one step")
	}
	now := time.Now().UTC()
	if wf.CreatedAt.IsZero() {
		wf.CreatedAt = now
	}
	wf.UpdatedAt = now
	wf.Name = name

	data, err := json.MarshalIndent(wf, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal workflow: %w", err)
	}
	return os.WriteFile(workflowPath(root, name), data, 0o600)
}

// LoadWorkflow reads a workflow definition from the store directory.
func LoadWorkflow(root, name string) (Workflow, error) {
	data, err := os.ReadFile(workflowPath(root, name))
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return Workflow{}, fmt.Errorf("workflow %q not found", name)
		}
		return Workflow{}, fmt.Errorf("read workflow: %w", err)
	}
	var wf Workflow
	if err := json.Unmarshal(data, &wf); err != nil {
		return Workflow{}, fmt.Errorf("unmarshal workflow: %w", err)
	}
	return wf, nil
}

// DeleteWorkflow removes a workflow definition from the store directory.
func DeleteWorkflow(root, name string) error {
	path := workflowPath(root, name)
	if err := os.Remove(path); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("workflow %q not found", name)
		}
		return fmt.Errorf("delete workflow: %w", err)
	}
	return nil
}

// ListWorkflows returns the names of all saved workflows in the store directory.
func ListWorkflows(root string) ([]string, error) {
	entries, err := os.ReadDir(root)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, nil
		}
		return nil, fmt.Errorf("list workflows: %w", err)
	}
	var names []string
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		const suffix = ".workflow.json"
		if n := e.Name(); len(n) > len(suffix) && n[len(n)-len(suffix):] == suffix {
			names = append(names, n[:len(n)-len(suffix)])
		}
	}
	return names, nil
}
