package cmd

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"

	"envport/internal/store"
)

func newMergeCmd(m Manager) *cobra.Command {
	var overwrite bool

	cmd := &cobra.Command{
		Use:   "merge <src> <dst>",
		Short: "Merge variables from src snapshot into dst snapshot",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			src, dst := args[0], args[1]

			strategy := store.MergeSkip
			if overwrite {
				strategy = store.MergeOverwrite
			}

			if err := m.Merge(dst, src, strategy); err != nil {
				if errors.Is(err, ErrNotFound) {
					return fmt.Errorf("snapshot not found")
				}
				return err
			}

			mode := "skipping"
			if overwrite {
				mode = "overwriting"
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Merged %q into %q (%s conflicts)\n", src, dst, mode)
			return nil
		},
	}

	cmd.Flags().BoolVarP(&overwrite, "overwrite", "o", false, "overwrite existing keys in dst with values from src")
	return cmd
}
