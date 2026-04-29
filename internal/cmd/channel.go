package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newChannelCmd(m Manager) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "channel",
		Short: "Manage release channel assignments for snapshots",
	}
	cmd.AddCommand(
		newChannelSetCmd(m),
		newChannelGetCmd(m),
		newChannelClearCmd(m),
		newChannelListCmd(m),
	)
	return cmd
}

func newChannelSetCmd(m Manager) *cobra.Command {
	return &cobra.Command{
		Use:   "set <name> <channel>",
		Short: "Set the release channel for a snapshot",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := m.SetChannel(args[0], args[1]); err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "channel %q set for snapshot %q\n", args[1], args[0])
			return nil
		},
	}
}

func newChannelGetCmd(m Manager) *cobra.Command {
	return &cobra.Command{
		Use:   "get <name>",
		Short: "Get the release channel for a snapshot",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ch, err := m.GetChannel(args[0])
			if err != nil {
				return err
			}
			if ch == "" {
				fmt.Fprintf(cmd.OutOrStdout(), "no channel set\n")
			} else {
				fmt.Fprintf(cmd.OutOrStdout(), "%s\n", ch)
			}
			return nil
		},
	}
}

func newChannelClearCmd(m Manager) *cobra.Command {
	return &cobra.Command{
		Use:   "clear <name>",
		Short: "Clear the release channel for a snapshot",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := m.ClearChannel(args[0]); err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "channel cleared for snapshot %q\n", args[0])
			return nil
		},
	}
}

func newChannelListCmd(m Manager) *cobra.Command {
	return &cobra.Command{
		Use:   "list <channel>",
		Short: "List snapshots assigned to a channel",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			names, err := m.ListByChannel(args[0])
			if err != nil {
				return err
			}
			if len(names) == 0 {
				fmt.Fprintf(cmd.OutOrStdout(), "no snapshots in channel %q\n", args[0])
				return nil
			}
			for _, n := range names {
				fmt.Fprintln(cmd.OutOrStdout(), n)
			}
			return nil
		},
	}
}
