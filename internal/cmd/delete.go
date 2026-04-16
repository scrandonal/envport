package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"envport/internal/store"
)

func newDeleteCmd(mgr *store.Manager) *cobra.Command {
	var force bool
	cmd := &cobra.Command{
		Use:     "delete <name>",
		Short:   "Delete a saved snapshot",
		Aliases: []string{"rm"},
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]
			if !force {
				fmt.Fprintf(cmd.OutOrStdout(), "Delete snapshot %q? [y/N]: ", name)
				var resp string
				fmt.Fscan(cmd.InOrStdin(), &resp)
				if resp != "y" && resp != "Y" {
					fmt.Fprintln(cmd.OutOrStdout(), "Aborted.")
					return nil
				}
			}
			if err := mgr.Delete(name); err != nil {
				return fmt.Errorf("delete: %w", err)
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Deleted snapshot %q\n", name)
			return nil
		},
	}
	cmd.Flags().BoolVarP(&force, "force", "f", false, "Skip confirmation prompt")
	return cmd
}
