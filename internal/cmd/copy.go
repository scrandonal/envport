package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newCopyCmd(m manager) *cobra.Command {
	var overwrite bool

	cmd := &cobra.Command{
		Use:   "copy <src> <dst>",
		Short: "Copy a snapshot to a new name",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			src, dst := args[0], args[1]

			snap, err := m.Load(src)
			if err != nil {
				return fmt.Errorf("copy: load %q: %w", src, err)
			}

			names, err := m.List()
			if err != nil {
				return fmt.Errorf("copy: list snapshots: %w", err)
			}

			if !overwrite {
				for _, n := range names {
					if n == dst {
						return fmt.Errorf("copy: %q already exists (use --overwrite to replace)", dst)
					}
				}
			}

			if err := m.Save(dst, snap); err != nil {
				return fmt.Errorf("copy: save %q: %w", dst, err)
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Copied %q → %q\n", src, dst)
			return nil
		},
	}

	cmd.Flags().BoolVar(&overwrite, "overwrite", false, "Overwrite destination if it exists")
	return cmd
}
