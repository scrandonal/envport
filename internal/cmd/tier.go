package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

func newTierCmd(m Manager) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tier",
		Short: "Manage snapshot tiers (free, standard, premium, enterprise)",
	}
	cmd.AddCommand(newTierSetCmd(m))
	cmd.AddCommand(newTierGetCmd(m))
	cmd.AddCommand(newTierClearCmd(m))
	cmd.AddCommand(newTierListCmd(m))
	return cmd
}

func newTierSetCmd(m Manager) *cobra.Command {
	return &cobra.Command{
		Use:   "set <snapshot> <tier>",
		Short: "Set the tier of a snapshot",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := m.SetTier(args[0], args[1]); err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "tier set to %q for snapshot %q\n", args[1], args[0])
			return nil
		},
	}
}

func newTierGetCmd(m Manager) *cobra.Command {
	return &cobra.Command{
		Use:   "get <snapshot>",
		Short: "Get the tier of a snapshot",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			tier, err := m.GetTier(args[0])
			if err != nil {
				return err
			}
			if tier == "" {
				fmt.Fprintf(cmd.OutOrStdout(), "no tier set\n")
			} else {
				fmt.Fprintf(cmd.OutOrStdout(), "%s\n", tier)
			}
			return nil
		},
	}
}

func newTierClearCmd(m Manager) *cobra.Command {
	return &cobra.Command{
		Use:   "clear <snapshot>",
		Short: "Clear the tier of a snapshot",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := m.ClearTier(args[0]); err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "tier cleared for snapshot %q\n", args[0])
			return nil
		},
	}
}

func newTierListCmd(m Manager) *cobra.Command {
	return &cobra.Command{
		Use:   "list <tier>",
		Short: "List snapshots with a given tier",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			names, err := m.ListByTier(args[0])
			if err != nil {
				return err
			}
			if len(names) == 0 {
				fmt.Fprintf(cmd.OutOrStdout(), "no snapshots with tier %q\n", args[0])
				return nil
			}
			fmt.Fprintln(cmd.OutOrStdout(), strings.Join(names, "\n"))
			return nil
		},
	}
}
