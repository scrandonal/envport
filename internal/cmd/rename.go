package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newRenameCmd(m managerIface) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rename <old> <new>",
		Short: "Rename a saved snapshot",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			oldName := args[0]
			newName := args[1]
			if err := m.Rename(oldName, newName); err != nil {
				return fmt.Errorf("rename: %w", err)
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Renamed %q to %q\n", oldName, newName)
			return nil
		},
	}
	return cmd
}
