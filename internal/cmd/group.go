package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

func newGroupCmd(m Manager) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "group",
		Short: "Manage snapshot groups",
	}
	cmd.AddCommand(newGroupCreateCmd(m))
	cmd.AddCommand(newGroupListCmd(m))
	cmd.AddCommand(newGroupDeleteCmd(m))
	cmd.AddCommand(newGroupShowCmd(m))
	return cmd
}

func newGroupCreateCmd(m Manager) *cobra.Command {
	return &cobra.Command{
		Use:   "create <group> <snap1> [snap2...]",
		Short: "Create a named group of snapshots",
		Args:  cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := m.CreateGroup(args[0], args[1:]); err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "group %q created with %d snapshot(s)\n", args[0], len(args[1:]))
			return nil
		},
	}
}

func newGroupListCmd(m Manager) *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List all groups",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			groups, err := m.ListGroups()
			if err != nil {
				return err
			}
			if len(groups) == 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "no groups found")
				return nil
			}
			for _, g := range groups {
				fmt.Fprintln(cmd.OutOrStdout(), g)
			}
			return nil
		},
	}
}

func newGroupDeleteCmd(m Manager) *cobra.Command {
	return &cobra.Command{
		Use:   "delete <group>",
		Short: "Delete a group",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := m.DeleteGroup(args[0]); err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "group %q deleted\n", args[0])
			return nil
		},
	}
}

func newGroupShowCmd(m Manager) *cobra.Command {
	return &cobra.Command{
		Use:   "show <group>",
		Short: "Show snapshots in a group",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			snaps, err := m.GetGroup(args[0])
			if err != nil {
				return err
			}
			fmt.Fprintln(cmd.OutOrStdout(), strings.Join(snaps, "\n"))
			return nil
		},
	}
}
