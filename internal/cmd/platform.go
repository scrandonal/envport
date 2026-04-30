package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newPlatformCmd(m Manager) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "platform",
		Short: "Manage platform tags for snapshots",
	}
	cmd.AddCommand(
		newPlatformSetCmd(m),
		newPlatformGetCmd(m),
		newPlatformClearCmd(m),
		newPlatformListCmd(m),
	)
	return cmd
}

func newPlatformSetCmd(m Manager) *cobra.Command {
	return &cobra.Command{
		Use:   "set <snapshot> <platform>",
		Short: "Set the platform for a snapshot",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := m.SetPlatform(args[0], args[1]); err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "platform set to %q for snapshot %q\n", args[1], args[0])
			return nil
		},
	}
}

func newPlatformGetCmd(m Manager) *cobra.Command {
	return &cobra.Command{
		Use:   "get <snapshot>",
		Short: "Get the platform for a snapshot",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			p, err := m.GetPlatform(args[0])
			if err != nil {
				return err
			}
			if p == "" {
				fmt.Fprintln(cmd.OutOrStdout(), "(no platform set)")
			} else {
				fmt.Fprintln(cmd.OutOrStdout(), p)
			}
			return nil
		},
	}
}

func newPlatformClearCmd(m Manager) *cobra.Command {
	return &cobra.Command{
		Use:   "clear <snapshot>",
		Short: "Clear the platform for a snapshot",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := m.ClearPlatform(args[0]); err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "platform cleared for snapshot %q\n", args[0])
			return nil
		},
	}
}

func newPlatformListCmd(m Manager) *cobra.Command {
	return &cobra.Command{
		Use:   "list <platform>",
		Short: "List snapshots by platform",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			names, err := m.ListByPlatform(args[0])
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
