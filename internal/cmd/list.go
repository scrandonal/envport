package cmd

import (
	"fmt"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"envport/internal/store"
)

func newListCmd(mgr *store.Manager) *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List all saved snapshots",
		Aliases: []string{"ls"},
		RunE: func(cmd *cobra.Command, args []string) error {
			names, err := mgr.List()
			if err != nil {
				return fmt.Errorf("list: %w", err)
			}
			if len(names) == 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "No snapshots saved.")
				return nil
			}
			w := tabwriter.NewWriter(cmd.OutOrStdout(), 0, 0, 2, ' ', 0)
			fmt.Fprintln(w, "NAME")
			for _, n := range names {
				fmt.Fprintln(w, n)
			}
			if err := w.Flush(); err != nil {
				return fmt.Errorf("list: flushing output: %w", err)
			}
			fmt.Fprintf(cmd.OutOrStdout(), "\n%d snapshot(s) total.\n", len(names))
			return nil
		},
	}
}
