package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newTagCmd(m Manager) *cobra.Command {
	var remove bool

	cmd := &cobra.Command{
		Use:   "tag <snapshot> <tag>",
		Short: "Add or remove a tag on a snapshot",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]
			tag := args[1]

			snap, err := m.Load(name)
			if err != nil {
				return fmt.Errorf("snapshot %q not found", name)
			}

			if remove {
				snap.RemoveTag(tag)
				fmt.Fprintf(cmd.OutOrStdout(), "removed tag %q from %q\n", tag, name)
			} else {
				snap.AddTag(tag)
				fmt.Fprintf(cmd.OutOrStdout(), "tagged %q with %q\n", name, tag)
			}

			if err := m.Save(name, snap); err != nil {
				return fmt.Errorf("failed to save snapshot: %w", err)
			}
			return nil
		},
	}

	cmd.Flags().BoolVarP(&remove, "remove", "r", false, "remove the tag instead of adding it")
	return cmd
}
