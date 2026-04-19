package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"envport/internal/store"
)

func newScheduleCmd(mgr Manager) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "schedule",
		Short: "Manage snapshot schedules",
	}
	cmd.AddCommand(newScheduleSetCmd(mgr))
	cmd.AddCommand(newScheduleGetCmd(mgr))
	cmd.AddCommand(newScheduleClearCmd(mgr))
	cmd.AddCommand(newScheduleListCmd(mgr))
	return cmd
}

func newScheduleSetCmd(mgr Manager) *cobra.Command {
	return &cobra.Command{
		Use:   "set <name> <cron> <action>",
		Short: "Set a schedule for a snapshot",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			s := store.Schedule{Cron: args[1], Action: args[2]}
			if err := mgr.SetSchedule(args[0], s); err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "schedule set for %s\n", args[0])
			return nil
		},
	}
}

func newScheduleGetCmd(mgr Manager) *cobra.Command {
	return &cobra.Command{
		Use:   "get <name>",
		Short: "Get schedule for a snapshot",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			s, err := mgr.GetSchedule(args[0])
			if err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "cron: %s\naction: %s\n",Action)
			return nil
		},
	}
}

func newScheduleClearCmd(mgr Manager) *cobra.Command {
	return &cobra.Command{
		Use:   "clear <name>",
		Short: "Clear schedule for a snapshot",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := mgr.ClearSchedule(args[0]); err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "schedule cleared for %s\n", args[0])
			return nil
		},
	}
}

func newScheduleListCmd(mgr Manager) *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List snapshots with schedules",
		RunE: func(cmd *cobra.Command, args []string) error {
			names, err := mgr.ListScheduled()
			if err != nil {
				return err
			}
			if len(names) == 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "no schedules set")
				return nil
			}
			for _, n := range names {
				fmt.Fprintln(cmd.OutOrStdout(), n)
			}
			return nil
		},
	}
}
