package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

func newStageCmd(m Manager) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "stage",
		Short: "Manage deployment stage metadata for snapshots",
	}
	cmd.AddCommand(
		newStageSetCmd(m),
		newStageGetCmd(m),
		newStageClearCmd(m),
		newStageListCmd(m),
	)
	return cmd
}

func newStageSetCmd(m Manager) *cobra.Command {
	return &cobra.Command{
		Use:   "set <name> <stage>",
		Short: "Set the deployment stage for a snapshot",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return m.SetStage(args[0], args[1])
		},
	}
}

func newStageGetCmd(m Manager) *cobra.Command {
	return &cobra.Command{
		Use:   "get <name>",
		Short: "Get the deployment stage for a snapshot",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			stage, err := m.GetStage(args[0])
			if err != nil {
				return err
			}
			if stage == "" {
				fmt.Fprintln(cmd.OutOrStdout(), "(no stage set)")
			} else {
				fmt.Fprintln(cmd.OutOrStdout(), stage)
			}
			return nil
		},
	}
}

func newStageClearCmd(m Manager) *cobra.Command {
	return &cobra.Command{
		Use:   "clear <name>",
		Short: "Clear the deployment stage for a snapshot",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return m.ClearStage(args[0])
		},
	}
}

func newStageListCmd(m Manager) *cobra.Command {
	return &cobra.Command{
		Use:   "list <stage>",
		Short: "List snapshots with a given deployment stage",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			names, err := m.ListByStage(args[0])
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
