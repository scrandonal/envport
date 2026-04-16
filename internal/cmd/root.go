package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"envport/internal/store"
)

func NewRootCmd() *cobra.Command {
	root := &cobra.Command{
		Use:   "envport",
		Short: "Snapshot and restore environment variable sets",
		SilenceUsage: true,
	}

	s, err := store.Default()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error initialising store: %v\n", err)
		os.Exit(1)
	}
	mgr := store.NewManager(s)

	root.AddCommand(
		newSaveCmd(mgr),
		newLoadCmd(mgr),
		newListCmd(mgr),
		newDeleteCmd(mgr),
	)
	return root
}
