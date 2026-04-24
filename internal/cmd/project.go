package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newProjectCmd(m Manager) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "project",
		Short: "Manage project associations for snapshots",
	}
	cmd.AddCommand(
		newProjectSetCmd(m),
		newProjectGetCmd(m),
		newProjectClearCmd(m),
		newProjectListCmd(m),
	)
	return cmd
}

func newProjectSetCmd(m Manager) *cobra.Command {
	return &cobra.Command{
		Use:   "set <snapshot> <project>",
		Short: "Assign a project to a snapshot",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := m.SetProject(args[0], args[1]); err != nil {
				return fmt.Errorf("set project: %w", err)
			}
			fmt.Fprintf(cmd.OutOrStdout(), "project %q assigned to %q\n", args[1], args[0])
			return nil
		},
	}
}

func newProjectGetCmd(m Manager) *cobra.Command {
	return &cobra.Command{
		Use:   "get <snapshot>",
		Short: "Show the project assigned to a snapshot",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			p, err := m.GetProject(args[0])
			if err != nil {
				return fmt.Errorf("get project: %w", err)
			}
			if p == "" {
				fmt.Fprintln(cmd.OutOrStdout(), "(none)")
			} else {
				fmt.Fprintln(cmd.OutOrStdout(), p)
			}
			return nil
		},
	}
}

func newProjectClearCmd(m Manager) *cobra.Command {
	return &cobra.Command{
		Use:   "clear <snapshot>",
		Short: "Remove the project association from a snapshot",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := m.ClearProject(args[0]); err != nil {
				return fmt.Errorf("clear project: %w", err)
			}
			fmt.Fprintf(cmd.OutOrStdout(), "project cleared for %q\n", args[0])
			return nil
		},
	}
}

func newProjectListCmd(m Manager) *cobra.Command {
	return &cobra.Command{
		Use:   "list <project>",
		Short: "List snapshots assigned to a project",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			names, err := m.ListByProject(args[0])
			if err != nil {
				return fmt.Errorf("list by project: %w", err)
			}
			if len(names) == 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "(none)")
				return nil
			}
			for _, n := range names {
				fmt.Fprintln(cmd.OutOrStdout(), n)
			}
			return nil
		},
	}
}
