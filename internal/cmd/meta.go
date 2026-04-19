package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newMetaCmd(m Manager) *cobra.Command {
	var clear bool

	cmd := &cobra.Command{
		Use:   "meta <name> [description]",
		Short: "Get or set a description for a snapshot",
		Args:  cobra.RangeArgs(1, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]

			if clear {
				return m.ClearMeta(name)
			}

			if len(args) == 2 {
				return m.SetMeta(name, args[1])
			}

			meta, err := m.GetMeta(name)
			if err != nil {
				return err
			}
			if meta.Description == "" {
				fmt.Fprintln(cmd.OutOrStdout(), "(no description set)")
			} else {
				fmt.Fprintln(cmd.OutOrStdout(), meta.Description)
			}
			if !meta.CreatedAt.IsZero() {
				fmt.Fprintf(cmd.OutOrStdout(), "Created: %s\n", meta.CreatedAt.Format("2006-01-02 15:04:05"))
				fmt.Fprintf(cmd.OutOrStdout(), "Updated: %s\n", meta.UpdatedAt.Format("2006-01-02 15:04:05"))
			}
			return nil
		},
	}

	cmd.Flags().BoolVar(&clear, "clear", false, "Clear the description")
	return cmd
}
