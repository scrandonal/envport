package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

func newLabelCmd(m Manager) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "label",
		Short: "Manage labels on snapshots",
	}
	cmd.AddCommand(newLabelAddCmd(m))
	cmd.AddCommand(newLabelRemoveCmd(m))
	cmd.AddCommand(newLabelListCmd(m))
	cmd.AddCommand(newLabelFindCmd(m))
	return cmd
}

func newLabelAddCmd(m Manager) *cobra.Command {
	return &cobra.Command{
		Use:   "add <snapshot> <label>",
		Short: "Add a label to a snapshot",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := m.AddLabel(args[0], args[1]); err != nil {
				return fmt.Errorf("add label: %w", err)
			}
			fmt.Fprintf(cmd.OutOrStdout(), "label %q added to %q\n", args[1], args[0])
			return nil
		},
	}
}

func newLabelRemoveCmd(m Manager) *cobra.Command {
	return &cobra.Command{
		Use:   "remove <snapshot> <label>",
		Short: "Remove a label from a snapshot",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := m.RemoveLabel(args[0], args[1]); err != nil {
				return fmt.Errorf("remove label: %w", err)
			}
			fmt.Fprintf(cmd.OutOrStdout(), "label %q removed from %q\n", args[1], args[0])
			return nil
		},
	}
}

func newLabelListCmd(m Manager) *cobra.Command {
	return &cobra.Command{
		Use:   "list <snapshot>",
		Short: "List labels on a snapshot",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			labels, err := m.GetLabels(args[0])
			if err != nil {
				return fmt.Errorf("get labels: %w", err)
			}
			if len(labels) == 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "no labels")
				return nil
			}
			fmt.Fprintln(cmd.OutOrStdout(), strings.Join(labels, "\n"))
			return nil
		},
	}
}

func newLabelFindCmd(m Manager) *cobra.Command {
	return &cobra.Command{
		Use:   "find <label>",
		Short: "Find snapshots with a given label",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			names, err := m.ListByLabel(args[0])
			if err != nil {
				return fmt.Errorf("find by label: %w", err)
			}
			if len(names) == 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "no snapshots found")
				return nil
			}
			fmt.Fprintln(cmd.OutOrStdout(), strings.Join(names, "\n"))
			return nil
		},
	}
}
