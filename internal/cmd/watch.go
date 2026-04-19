package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newWatchCmd(m Manager) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "watch",
		Short: "Manage snapshot watch markers",
	}

	setCmd := &cobra.Command{
		Use:   "set <name>",
		Short: "Mark a snapshot as watched",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := m.SetWatch(args[0]); err != nil {
				return fmt.Errorf("watch set: %w", err)
			}
			fmt.Fprintf(cmd.OutOrStdout(), "watching %s\n", args[0])
			return nil
		},
	}

	clearCmd := &cobra.Command{
		Use:   "clear <name>",
		Short: "Remove watch marker from a snapshot",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := m.ClearWatch(args[0]); err != nil {
				return fmt.Errorf("watch clear: %w", err)
			}
			fmt.Fprintf(cmd.OutOrStdout(), "cleared watch for %s\n", args[0])
			return nil
		},
	}

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List all watched snapshots",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			events, err := m.ListWatched()
			if err != nil {
				return fmt.Errorf("watch list: %w", err)
			}
			if len(events) == 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "no watched snapshots")
				return nil
			}
			for _, e := range events {
				fmt.Fprintf(cmd.OutOrStdout(), "%s\t%s\n", e.Name, e.ChangedAt.Format("2006-01-02 15:04:05"))
			}
			return nil
		},
	}

	cmd.AddCommand(setCmd, clearCmd, listCmd)
	return cmd
}
