package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/user/envport/internal/snapshot"
)

func newEnvDiffCmd(m Manager) *cobra.Command {
	var onlyChanged bool

	cmd := &cobra.Command{
		Use:   "envdiff <name>",
		Short: "Compare a snapshot against the current environment",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]
			environ := snapshot.OSEnviron()
			d, err := m.DiffWithEnviron(name, environ)
			if err != nil {
				return fmt.Errorf("envdiff: %w", err)
			}
			if !d.HasChanges() {
				fmt.Fprintln(cmd.OutOrStdout(), "no differences")
				return nil
			}
			lines := d.Summary()
			if onlyChanged {
				lines = filterChanged(lines)
			}
			for _, l := range lines {
				fmt.Fprintln(cmd.OutOrStdout(), l)
			}
			return nil
		},
	}
	cmd.Flags().BoolVar(&onlyChanged, "changed", false, "show only changed keys")
	_ = os.Stderr // suppress unused import
	return cmd
}

func filterChanged(lines []string) []string {
	var out []string
	for _, l := range lines {
		if len(l) > 0 && l[0] == '~' {
			out = append(out, l)
		}
	}
	return out
}
