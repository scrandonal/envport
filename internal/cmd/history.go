package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	"github.com/spf13/cobra"
)

type historyStore interface {
	ReadHistory() ([]historyEntry, error)
	ClearHistory() error
}

type historyEntry interface {
	GetName() string
	GetOperation() string
	GetTimestamp() time.Time
}

func newHistoryCmd(m Manager) *cobra.Command {
	var clear bool

	cmd := &cobra.Command{
		Use:   "history",
		Short: "Show operation history for snapshots",
		RunE: func(cmd *cobra.Command, args []string) error {
			if clear {
				if err := m.ClearHistory(); err != nil {
					return fmt.Errorf("clear history: %w", err)
				}
				fmt.Fprintln(cmd.OutOrStdout(), "History cleared.")
				return nil
			}

			entries, err := m.ListHistory()
			if err != nil {
				return fmt.Errorf("read history: %w", err)
			}
			if len(entries) == 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "No history recorded.")
				return nil
			}

			w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
			fmt.Fprintln(w, "TIMESTAMP\tOPERATION\tNAME")
			for _, e := range entries {
				fmt.Fprintf(w, "%s\t%s\t%s\n",
					e.Timestamp.Format(time.RFC3339),
					e.Operation,
					e.Name,
				)
			}
			return w.Flush()
		},
	}

	cmd.Flags().BoolVar(&clear, "clear", false, "Clear all history entries")
	return cmd
}
