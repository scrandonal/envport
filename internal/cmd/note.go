package cmd

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"

	"envport/internal/store"
)

func newNoteCmd(m Manager) *cobra.Command {
	var clear bool

	cmd := &cobra.Command{
		Use:   "note <name> [text]",
		Short: "Get or set a note on a snapshot",
		Args:  cobra.RangeArgs(1, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]

			if clear {
				return m.ClearNote(name)
			}

			if len(args) == 2 {
				return m.SetNote(name, args[1])
			}

			// read mode
			note, err := m.GetNote(name)
			if errors.Is(err, store.ErrNoteNotFound) {
				fmt.Fprintln(cmd.OutOrStdout(), "(no note)")
				return nil
			}
			if err != nil {
				return err
			}
			fmt.Fprintln(cmd.OutOrStdout(), note)
			return nil
		},
	}

	cmd.Flags().BoolVar(&clear, "clear", false, "Remove the note from the snapshot")
	return cmd
}
