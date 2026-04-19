package cmd

import (
	"fmt"
	"sort"

	"github.com/spf13/cobra"
)

func newCompareCmd(m Manager) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "compare <src> <dst>",
		Short: "Show a detailed diff between two snapshots",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			d, err := m.Compare(args[0], args[1])
			if err != nil {
				return err
			}

			if len(d.Added) == 0 && len(d.Removed) == 0 && len(d.Changed) == 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "No differences found.")
				return nil
			}

			keys := func(m map[string]string) []string {
				out := make([]string, 0, len(m))
				for k := range m {
					out = append(out, k)
				}
				sort.Strings(out)
				return out
			}

			w := cmd.OutOrStdout()
			for _, k := range keys(d.Added) {
				fmt.Fprintf(w, "+ %s=%s\n", k, d.Added[k])
			}
			for _, k := range keys(d.Removed) {
				fmt.Fprintf(w, "- %s=%s\n", k, d.Removed[k])
			}
			changedKeys := make([]string, 0, len(d.Changed))
			for k := range d.Changed {
				changedKeys = append(changedKeys, k)
			}
			sort.Strings(changedKeys)
			for _, k := range changedKeys {
				c := d.Changed[k]
				fmt.Fprintf(w, "~ %s: %s -> %s\n", k, c.Old, c.New)
			}
			return nil
		},
	}
	return cmd
}
