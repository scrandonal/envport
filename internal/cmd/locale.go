package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

func newLocaleCmd(mgr Manager) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "locale",
		Short: "Manage locale metadata for snapshots",
	}
	cmd.AddCommand(newLocaleSetCmd(mgr))
	cmd.AddCommand(newLocaleGetCmd(mgr))
	cmd.AddCommand(newLocaleClearCmd(mgr))
	cmd.AddCommand(newLocaleListCmd(mgr))
	return cmd
}

func newLocaleSetCmd(mgr Manager) *cobra.Command {
	return &cobra.Command{
		Use:   "set <snapshot> <locale>",
		Short: "Set the locale for a snapshot",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := mgr.SetLocale(args[0], args[1]); err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "locale set to %q for snapshot %q\n", args[1], args[0])
			return nil
		},
	}
}

func newLocaleGetCmd(mgr Manager) *cobra.Command {
	return &cobra.Command{
		Use:   "get <snapshot>",
		Short: "Get the locale for a snapshot",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			locale, err := mgr.GetLocale(args[0])
			if err != nil {
				return err
			}
			if locale == "" {
				fmt.Fprintf(cmd.OutOrStdout(), "no locale set for snapshot %q\n", args[0])
			} else {
				fmt.Fprintln(cmd.OutOrStdout(), locale)
			}
			return nil
		},
	}
}

func newLocaleClearCmd(mgr Manager) *cobra.Command {
	return &cobra.Command{
		Use:   "clear <snapshot>",
		Short: "Clear the locale for a snapshot",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := mgr.ClearLocale(args[0]); err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "locale cleared for snapshot %q\n", args[0])
			return nil
		},
	}
}

func newLocaleListCmd(mgr Manager) *cobra.Command {
	return &cobra.Command{
		Use:   "list <locale>",
		Short: "List snapshots with a given locale",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			names, err := mgr.ListByLocale(args[0])
			if err != nil {
				return err
			}
			if len(names) == 0 {
				fmt.Fprintf(cmd.OutOrStdout(), "no snapshots with locale %q\n", args[0])
				return nil
			}
			fmt.Fprintln(cmd.OutOrStdout(), strings.Join(names, "\n"))
			return nil
		},
	}
}
