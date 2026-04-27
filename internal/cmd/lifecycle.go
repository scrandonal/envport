package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newLifecycleCmd(m Manager) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "lifecycle",
		Short: "Manage snapshot lifecycle stages",
	}
	cmd.AddCommand(
		newLifecycleSetCmd(m),
		newLifecycleGetCmd(m),
		newLifecycleClearCmd(m),
		newLifecycleListCmd(m),
	)
	return cmd
}

func newLifecycleSetCmd(m Manager) *cobra.Command {
	return &cobra.Command{
		Use:   "set <name> <stage>",
		Short: "Set lifecycle stage (draft|active|deprecated|archived|retired)",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := m.SetLifecycle(args[0], args[1]); err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "lifecycle of %q set to %q\n", args[0], args[1])
			return nil
		},
	}
}

func newLifecycleGetCmd(m Manager) *cobra.Command {
	return &cobra.Command{
		Use:   "get <name>",
		Short: "Get the lifecycle stage of a snapshot",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			lc, err := m.GetLifecycle(args[0])
			if err != nil {
				return err
			}
			if lc == "" {
				fmt.Fprintln(cmd.OutOrStdout(), "(no lifecycle stage set)")
			} else {
				fmt.Fprintln(cmd.OutOrStdout(), lc)
			}
			return nil
		},
	}
}

func newLifecycleClearCmd(m Manager) *cobra.Command {
	return &cobra.Command{
		Use:   "clear <name>",
		Short: "Clear the lifecycle stage of a snapshot",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := m.ClearLifecycle(args[0]); err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "lifecycle cleared for %q\n", args[0])
			return nil
		},
	}
}

func newLifecycleListCmd(m Manager) *cobra.Command {
	return &cobra.Command{
		Use:   "list <stage>",
		Short: "List snapshots with a given lifecycle stage",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			names, err := m.ListByLifecycle(args[0])
			if err != nil {
				return err
			}
			if len(names) == 0 {
				fmt.Fprintf(cmd.OutOrStdout(), "no snapshots with lifecycle stage %q\n", args[0])
				return nil
			}
			for _, n := range names {
				fmt.Fprintln(cmd.OutOrStdout(), n)
			}
			return nil
		},
	}
}
