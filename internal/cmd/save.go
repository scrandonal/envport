package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"envport/internal/snapshot"
	"envport/internal/store"
)

func newSaveCmd(mgr *store.Manager) *cobra.Command {
	return &cobra.Command{
		Use:   "save <name>",
		Short: "Snapshot current environment variables and save under <name>",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]
			snap := snapshot.FromEnviron(os.Environ())
			if err := mgr.Save(name, snap); err != nil {
				return fmt.Errorf("save: %w", err)
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Saved snapshot %q (%d vars)\n", name, len(snap.Vars))
			return nil
		},
	}
}
