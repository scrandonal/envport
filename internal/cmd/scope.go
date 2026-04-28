package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

func newScopeCmd(m Manager) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "scope",
		Short: "Manage snapshot scopes",
	}
	cmd.AddCommand(
		newScopeSetCmd(m),
		newScopeGetCmd(m),
		newScopeClearCmd(m),
		newScopeListCmd(m),
	)
	return cmd
}

func newScopeSetCmd(m Manager) *cobra.Command {
	return &cobra.Command{
		Use:   "set <name> <scope>",
		Short: "Set scope for a snapshot (global|local|session|user)",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := m.SetScope(args[0], args[1]); err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "scope set to %q for snapshot %q\n", args[1], args[0])
			return nil
		},
	}
}

func newScopeGetCmd(m Manager) *cobra.Command {
	return &cobra.Command{
		Use:   "get <name>",
		Short: "Get scope of a snapshot",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			scope, err := m.GetScope(args[0])
			if err != nil {
				return err
			}
			if scope == "" {
				fmt.Fprintln(cmd.OutOrStdout(), "(no scope set)")
			} else {
				fmt.Fprintln(cmd.OutOrStdout(), scope)
			}
			return nil
		},
	}
}

func newScopeClearCmd(m Manager) *cobra.Command {
	return &cobra.Command{
		Use:   "clear <name>",
		Short: "Clear scope of a snapshot",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := m.ClearScope(args[0]); err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "scope cleared for snapshot %q\n", args[0])
			return nil
		},
	}
}

func newScopeListCmd(m Manager) *cobra.Command {
	return &cobra.Command{
		Use:   "list <scope>",
		Short: "List snapshots with a given scope",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			names, err := m.ListByScope(args[0])
			if err != nil {
				return err
			}
			if len(names) == 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "(none)")
				return nil
			}
			fmt.Fprintln(cmd.OutOrStdout(), strings.Join(names, "\n"))
			return nil
		},
	}
}
