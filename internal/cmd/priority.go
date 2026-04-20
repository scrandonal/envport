package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

func newPriorityCmd(m Manager) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "priority",
		Short: "Manage snapshot priority levels",
	}
	cmd.AddCommand(
		newPrioritySetCmd(m),
		newPriorityGetCmd(m),
		newPriorityClearCmd(m),
		newPriorityListCmd(m),
	)
	return cmd
}

func newPrioritySetCmd(m Manager) *cobra.Command {
	return &cobra.Command{
		Use:   "set <name> <level>",
		Short: "Set priority level for a snapshot (low|normal|high|critical)",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := m.SetPriority(args[0], args[1]); err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "priority for %q set to %q\n", args[0], args[1])
			return nil
		},
	}
}

func newPriorityGetCmd(m Manager) *cobra.Command {
	return &cobra.Command{
		Use:   "get <name>",
		Short: "Get priority level for a snapshot",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			level, err := m.GetPriority(args[0])
			if err != nil {
				return err
			}
			fmt.Fprintln(cmd.OutOrStdout(), level)
			return nil
		},
	}
}

func newPriorityClearCmd(m Manager) *cobra.Command {
	return &cobra.Command{
		Use:   "clear <name>",
		Short: "Clear priority level for a snapshot",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := m.ClearPriority(args[0]); err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "priority cleared for %q\n", args[0])
			return nil
		},
	}
}

func newPriorityListCmd(m Manager) *cobra.Command {
	return &cobra.Command{
		Use:   "list <level>",
		Short: "List snapshots with a given priority level",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			names, err := m.ListByPriority(args[0])
			if err != nil {
				return err
			}
			if len(names) == 0 {
				fmt.Fprintf(cmd.OutOrStdout(), "no snapshots with priority %q\n", args[0])
				return nil
			}
			fmt.Fprintln(cmd.OutOrStdout(), strings.Join(names, "\n"))
			return nil
		},
	}
}
