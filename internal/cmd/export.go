package cmd

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/spf13/cobra"
)

func newExportCmd(m Manager) *cobra.Command {
	var format string

	cmd := &cobra.Command{
		Use:   "export <name>",
		Short: "Export a snapshot as shell export statements or dotenv format",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			snap, err := m.Load(args[0])
			if err != nil {
				return fmt.Errorf("snapshot %q not found", args[0])
			}

			keys := make([]string, 0, len(snap.Vars))
			for k := range snap.Vars {
				keys = append(keys, k)
			}
			sort.Strings(keys)

			w := cmd.OutOrStdout()
			for _, k := range keys {
				v := snap.Vars[k]
				switch strings.ToLower(format) {
				case "dotenv":
					fmt.Fprintf(w, "%s=%s\n", k, v)
				default:
					fmt.Fprintf(w, "export %s=%q\n", k, v)
				}
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&format, "format", "f", "shell", "Output format: shell or dotenv")
	_ = os.Stderr // suppress unused import
	return cmd
}
