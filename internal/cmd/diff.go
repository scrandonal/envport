package cmd

import (
	"fmt"
	"sort"

	"github.com/spf13/cobra"
)

func newDiffCmd(m Manager) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "diff <snapshot1> <snapshot2>",
		Short: "Show differences between two snapshots",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			src, err := m.Load(args[0])
			if err != nil {
				return fmt.Errorf("loading %q: %w", args[0], err)
			}
			dst, err := m.Load(args[1])
			if err != nil {
				return fmt.Errorf("loading %q: %w", args[1], err)
			}

			keys := unionKeys(src.Vars, dst.Vars)
			sort.Strings(keys)

			changes := 0
			for _, k := range keys {
				a, inA := src.Vars[k]
				b, inB := dst.Vars[k]
				switch {
				case inA && !inB:
					fmt.Fprintf(cmd.OutOrStdout(), "- %s=%s\n", k, a)
					changes++
				case !inA && inB:
					fmt.Fprintf(cmd.OutOrStdout(), "+ %s=%s\n", k, b)
					changes++
				case a != b:
					fmt.Fprintf(cmd.OutOrStdout(), "~ %s: %s -> %s\n", k, a, b)
					changes++
				}
			}
			if changes == 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "no differences")
			}
			return nil
		},
	}
	return cmd
}

func unionKeys(a, b map[string]string) []string {
	seen := make(map[string]struct{}, len(a)+len(b))
	for k := range a {
		seen[k] = struct{}{}
	}
	for k := range b {
		seen[k] = struct{}{}
	}
	keys := make([]string, 0, len(seen))
	for k := range seen {
		keys = append(keys, k)
	}
	return keys
}
