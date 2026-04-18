package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newCloneCmd(m Manager) *cobra.Command {
	var overwrite bool

	cmd := &cobra.Command{
		Use:   "clone <src> <dst>",
		Short: "Clone a snapshot to a new name",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			src, dst := args[0], args[1]
			if err := m.Clone(src, dst, overwrite); err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "cloned %q → %q\n", src, dst)
			return nil
		},
	}

	cmd.Flags().BoolVar(&overwrite, "overwrite", false, "overwrite destination if it exists")
	return cmd
}
