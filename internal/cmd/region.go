package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

func newRegionCmd(m Manager) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "region",
		Short: "Manage region metadata for snapshots",
	}
	cmd.AddCommand(
		newRegionSetCmd(m),
		newRegionGetCmd(m),
		newRegionClearCmd(m),
		newRegionListCmd(m),
	)
	return cmd
}

func newRegionSetCmd(m Manager) *cobra.Command {
	return &cobra.Command{
		Use:   "set <name> <region>",
		Short: "Set the region for a snapshot",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := m.SetRegion(args[0], args[1]); err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "region %q set for snapshot %q\n", args[1], args[0])
			return nil
		},
	}
}

func newRegionGetCmd(m Manager) *cobra.Command {
	return &cobra.Command{
		Use:   "get <name>",
		Short: "Get the region for a snapshot",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			r, err := m.GetRegion(args[0])
			if err != nil {
				return err
			}
			if r == "" {
				fmt.Fprintf(cmd.OutOrStdout(), "no region set for %q\n", args[0])
			} else {
				fmt.Fprintln(cmd.OutOrStdout(), r)
			}
			return nil
		},
	}
}

func newRegionClearCmd(m Manager) *cobra.Command {
	return &cobra.Command{
		Use:   "clear <name>",
		Short: "Clear the region for a snapshot",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := m.ClearRegion(args[0]); err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "region cleared for %q\n", args[0])
			return nil
		},
	}
}

func newRegionListCmd(m Manager) *cobra.Command {
	return &cobra.Command{
		Use:   "list <region>",
		Short: "List snapshots in a given region",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			names, err := m.ListByRegion(args[0])
			if err != nil {
				return err
			}
			if len(names) == 0 {
				fmt.Fprintf(cmd.OutOrStdout(), "no snapshots in region %q\n", args[0])
				return nil
			}
			fmt.Fprintln(cmd.OutOrStdout(), strings.Join(names, "\n"))
			return nil
		},
	}
}
