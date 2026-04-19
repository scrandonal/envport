package cmd

import (
	"fmt"
	"sort"

	"github.com/spf13/cobra"
)

func newStatsCmd(m Manager) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "stats <name>",
		Short: "Show statistics for a snapshot",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]
			snap, err := m.Load(name)
			if err != nil {
				return fmt.Errorf("snapshot %q not found", name)
			}

			type kv struct {
				key  string
				size int
			}

			var pairs []kv
			var totalBytes int
			for k, v := range snap.Vars {
				sz := len(k) + len(v)
				totalBytes += sz
				pairs = append(pairs, kv{k, sz})
			}
			sort.Slice(pairs, func(i, j int) bool {
				return pairs[i].key < pairs[j].key
			})

			fmt.Fprintf(cmd.OutOrStdout(), "Snapshot : %s\n", name)
			fmt.Fprintf(cmd.OutOrStdout(), "Keys     : %d\n", len(snap.Vars))
			fmt.Fprintf(cmd.OutOrStdout(), "Size     : %d bytes\n", totalBytes)

			verbose, _ := cmd.Flags().GetBool("verbose")
			if verbose {
				fmt.Fprintln(cmd.OutOrStdout(), "\nKey sizes:")
				for _, p := range pairs {
					fmt.Fprintf(cmd.OutOrStdout(), "  %-30s %d bytes\n", p.key, p.size)
				}
			}
			return nil
		},
	}
	cmd.Flags().BoolP("verbose", "v", false, "Show per-key sizes")
	return cmd
}
