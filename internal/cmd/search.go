package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

func newSearchCmd(m Manager) *cobra.Command {
	var matchAll bool

	cmd := &cobra.Command{
		Use:   "search <key=value|key> [key=value|key]...",
		Short: "Search snapshots by environment variable key or key=value",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			names, err := m.List()
			if err != nil {
				return err
			}

			for _, name := range names {
				snap, err := m.Load(name)
				if err != nil {
					continue
				}

				if matchesQuery(snap.Vars, args, matchAll) {
					fmt.Fprintln(cmd.OutOrStdout(), name)
				}
			}
			return nil
		},
	}

	cmd.Flags().BoolVarP(&matchAll, "all", "a", false, "require all terms to match (AND); default is OR")
	return cmd
}

func matchesQuery(vars map[string]string, terms []string, matchAll bool) bool {
	for _, term := range terms {
		matched := false
		if strings.Contains(term, "=") {
			parts := strings.SplitN(term, "=", 2)
			v, ok := vars[parts[0]]
			matched = ok && v == parts[1]
		} else {
			_, matched = vars[term]
		}
		if matchAll && !matched {
			return false
		}
		if !matchAll && matched {
			return true
		}
	}
	return matchAll
}
