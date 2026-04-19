package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	"github.com/spf13/cobra"
)

func newAuditCmd(m Manager) *cobra.Command {
	var clear bool
	cmd := &cobra.Command{
		Use:   "audit",
		Short: "Show or clear the audit log of snapshot operations",
		RunE: func(cmd *cobra.Command, args []string) error {
			if clear {
				if err := m.ClearAuditLog(); err != nil {
					return fmt.Errorf("clear audit log: %w", err)
				}
				fmt.Fprintln(cmd.OutOrStdout(), "audit log cleared")
				return nil
			}
			entries, err := m.AuditLog()
			if err != nil {
				return fmt.Errorf("read audit log: %w", err)
			}
			if len(entries) == 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "no audit entries")
				return nil
			}
			w := tabwriter.NewWriter(cmd.OutOrStdout(), 0, 0, 2, ' ', 0)
			fmt.Fprintln(w, "TIME\tOPERATION\tNAME\tDETAIL")
			for _, e := range entries {
				fmt.Fprintf(w, "%s\t%s\t%s\t%s\n",
					e.Time.Format(time.RFC3339), e.Operation, e.Name, e.Detail)
			}
			return w.Flush()
		},
	}
	cmd.Flags().BoolVar(&clear, "clear", false, "clear the audit log")
	_ = os.Stderr
	return cmd
}
