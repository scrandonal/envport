package cmd

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
)

func newRollbackCmd(m Manager) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rollback <name> [steps]",
		Short: "Restore a snapshot to a previous version from its history",
		Args:  cobra.RangeArgs(1, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]
			steps := 1
			if len(args) == 2 {
				var err error
				steps, err = strconv.Atoi(args[1])
				if err != nil || steps <= 0 {
					return fmt.Errorf("steps must be a positive integer")
				}
			}
			if err := m.Rollback(name, steps); err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Rolled back %q by %d step(s)\n", name, steps)
			return nil
		},
	}
	return cmd
}
