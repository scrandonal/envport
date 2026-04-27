package cmd

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
)

func newRetentionCmd(m Manager) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "retention",
		Short: "Manage retention policies for snapshots",
	}
	cmd.AddCommand(
		newRetentionSetCmd(m),
		newRetentionGetCmd(m),
		newRetentionClearCmd(m),
	)
	return cmd
}

func newRetentionSetCmd(m Manager) *cobra.Command {
	return &cobra.Command{
		Use:   "set <name> <days>",
		Short: "Set a retention policy (in days) for a snapshot",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			days, err := strconv.Atoi(args[1])
			if err != nil || days <= 0 {
				return fmt.Errorf("days must be a positive integer")
			}
			if err := m.SetRetention(args[0], days); err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "retention set to %d days for %q\n", days, args[0])
			return nil
		},
	}
}

func newRetentionGetCmd(m Manager) *cobra.Command {
	return &cobra.Command{
		Use:   "get <name>",
		Short: "Show the retention policy for a snapshot",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			p, err := m.GetRetention(args[0])
			if err != nil {
				return err
			}
			if p == nil {
				fmt.Fprintf(cmd.OutOrStdout(), "no retention policy set for %q\n", args[0])
				return nil
			}
			fmt.Fprintf(cmd.OutOrStdout(), "retention: %d days (set at %s)\n", p.Days, p.SetAt.Format("2006-01-02"))
			return nil
		},
	}
}

func newRetentionClearCmd(m Manager) *cobra.Command {
	return &cobra.Command{
		Use:   "clear <name>",
		Short: "Remove the retention policy from a snapshot",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := m.ClearRetention(args[0]); err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "retention policy cleared for %q\n", args[0])
			return nil
		},
	}
}
