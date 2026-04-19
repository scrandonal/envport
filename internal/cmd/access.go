package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/nicholasgasior/envport/internal/store"
)

func newAccessCmd(m Manager) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "access",
		Short: "Show or clear access records for a snapshot",
	}
	cmd.AddCommand(newAccessShowCmd(m))
	cmd.AddCommand(newAccessClearCmd(m))
	return cmd
}

func newAccessShowCmd(m Manager) *cobra.Command {
	return &cobra.Command{
		Use:   "show <name>",
		Short: "Show access record for a snapshot",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]
			s, ok := m.(*store.Manager)
			if !ok {
				return fmt.Errorf("access command requires store.Manager")
			}
			rec, err := store.GetAccess(s.Root(), name)
			if err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "load_count:  %d\n", rec.LoadCount)
			fmt.Fprintf(cmd.OutOrStdout(), "save_count:  %d\n", rec.SaveCount)
			if rec.LastLoaded != nil {
				fmt.Fprintf(cmd.OutOrStdout(), "last_loaded: %s\n", rec.LastLoaded.Format("2006-01-02 15:04:05"))
			}
			if rec.LastSaved != nil {
				fmt.Fprintf(cmd.OutOrStdout(), "last_saved:  %s\n", rec.LastSaved.Format("2006-01-02 15:04:05"))
			}
			return nil
		},
	}
}

func newAccessClearCmd(m Manager) *cobra.Command {
	return &cobra.Command{
		Use:   "clear <name>",
		Short: "Clear access record for a snapshot",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]
			s, ok := m.(*store.Manager)
			if !ok {
				return fmt.Errorf("access command requires store.Manager")
			}
			if err := store.ClearAccess(s.Root(), name); err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "access record cleared for %q\n", name)
			return nil
		},
	}
}
