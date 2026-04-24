package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newStatusCmd(m Manager) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "status",
		Short: "Manage snapshot status (active, deprecated, draft, archived)",
	}
	cmd.AddCommand(
		newStatusSetCmd(m),
		newStatusGetCmd(m),
		newStatusClearCmd(m),
		newStatusListCmd(m),
	)
	return cmd
}

func newStatusSetCmd(m Manager) *cobra.Command {
	return &cobra.Command{
		Use:   "set <name> <status>",
		Short: "Set the status of a snapshot",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := m.SetStatus(args[0], args[1]); err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "status of %q set to %q\n", args[0], args[1])
			return nil
		},
	}
}

func newStatusGetCmd(m Manager) *cobra.Command {
	return &cobra.Command{
		Use:   "get <name>",
		Short: "Get the status of a snapshot",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			s, err := m.GetStatus(args[0])
			if err != nil {
				return err
			}
			if s == "" {
				fmt.Fprintln(cmd.OutOrStdout(), "(no status set)")
			} else {
				fmt.Fprintln(cmd.OutOrStdout(), s)
			}
			return nil
		},
	}
}

func newStatusClearCmd(m Manager) *cobra.Command {
	return &cobra.Command{
		Use:   "clear <name>",
		Short: "Clear the status of a snapshot",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := m.ClearStatus(args[0]); err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "status cleared for %q\n", args[0])
			return nil
		},
	}
}

func newStatusListCmd(m Manager) *cobra.Command {
	return &cobra.Command{
		Use:   "list <status>",
		Short: "List snapshots with a given status",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			names, err := m.ListByStatus(args[0])
			if err != nil {
				return err
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
