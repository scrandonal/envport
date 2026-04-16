package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"envport/internal/store"
)

func newLoadCmd(mgr *store.Manager) *cobra.Command {
	var shellFmt string
	cmd := &cobra.Command{
		Use:   "load <name>",
		Short: "Print shell export statements for a saved snapshot",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]
			snap, err := mgr.Load(name)
			if err != nil {
				return fmt.Errorf("load: %w", err)
			}
			out := cmd.OutOrStdout()
			for k, v := range snap.Vars {
				switch shellFmt {
				case "fish":
					fmt.Fprintf(out, "set -x %s %q\n", k, v)
				default:
					fmt.Fprintf(out, "export %s=%q\n", k, v)
				}
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&shellFmt, "shell", "bash", "Shell format: bash, fish")
	return cmd
}
