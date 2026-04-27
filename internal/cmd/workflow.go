package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/user/envport/internal/store"
)

func newWorkflowCmd(m Manager) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "workflow",
		Short: "Manage workflows attached to snapshots",
	}
	cmd.AddCommand(
		newWorkflowSaveCmd(m),
		newWorkflowShowCmd(m),
		newWorkflowDeleteCmd(m),
		newWorkflowListCmd(m),
	)
	return cmd
}

func newWorkflowSaveCmd(m Manager) *cobra.Command {
	return &cobra.Command{
		Use:   "save <snapshot> <name> <step1,step2,...>",
		Short: "Save a named workflow for a snapshot",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			snap, name, raw := args[0], args[1], args[2]
			steps := strings.Split(raw, ",")
			wf := store.Workflow{Name: name, Steps: steps}
			if err := m.SaveWorkflow(snap, wf); err != nil {
				return fmt.Errorf("save workflow: %w", err)
			}
			fmt.Fprintf(cmd.OutOrStdout(), "workflow %q saved for snapshot %q\n", name, snap)
			return nil
		},
	}
}

func newWorkflowShowCmd(m Manager) *cobra.Command {
	return &cobra.Command{
		Use:   "show <snapshot> <name>",
		Short: "Show steps of a workflow",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			snap, name := args[0], args[1]
			wf, err := m.LoadWorkflow(snap, name)
			if err != nil {
				return fmt.Errorf("load workflow: %w", err)
			}
			fmt.Fprintf(cmd.OutOrStdout(), "workflow: %s\n", wf.Name)
			for i, s := range wf.Steps {
				fmt.Fprintf(cmd.OutOrStdout(), "  %d. %s\n", i+1, s)
			}
			return nil
		},
	}
}

func newWorkflowDeleteCmd(m Manager) *cobra.Command {
	return &cobra.Command{
		Use:   "delete <snapshot> <name>",
		Short: "Delete a workflow from a snapshot",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			snap, name := args[0], args[1]
			if err := m.DeleteWorkflow(snap, name); err != nil {
				return fmt.Errorf("delete workflow: %w", err)
			}
			fmt.Fprintf(cmd.OutOrStdout(), "workflow %q deleted from snapshot %q\n", name, snap)
			return nil
		},
	}
}

func newWorkflowListCmd(m Manager) *cobra.Command {
	return &cobra.Command{
		Use:   "list <snapshot>",
		Short: "List all workflows for a snapshot",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			snap := args[0]
			list, err := m.ListWorkflows(snap)
			if err != nil {
				return fmt.Errorf("list workflows: %w", err)
			}
			if len(list) == 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "no workflows")
				return nil
			}
			for _, wf := range list {
				fmt.Fprintf(cmd.OutOrStdout(), "%s (%d steps)\n", wf.Name, len(wf.Steps))
			}
			return nil
		},
	}
}
