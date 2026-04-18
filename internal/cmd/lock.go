package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newLockCmd(m Manager) *cobra.Command {
	var unlock bool

	cmd := &cobra.Command{
		Use:   "lock",
		Short: "Show or release the store lock",
		Long:  `Display whether the store is currently locked, or forcibly remove a stale lock.`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if unlock {
				if err := m.ForceUnlock(); err != nil {
					return fmt.Errorf("unlock: %w", err)
				}
				fmt.Fprintln(cmd.OutOrStdout(), "lock released")
				return nil
			}
			locked, err := m.IsLocked()
			if err != nil {
				return err
			}
			if locked {
				fmt.Fprintln(cmd.OutOrStdout(), "store is locked")
			} else {
				fmt.Fprintln(cmd.OutOrStdout(), "store is not locked")
			}
			return nil
		},
	}
	cmd.Flags().BoolVarP(&unlock, "unlock", "u", false, "forcibly remove a stale lock")
	return cmd
}
