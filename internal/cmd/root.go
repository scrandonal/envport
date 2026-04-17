package cmd

import (
	"github.com/nicholasgasior/envport/internal/store"
	"github.com/spf13/cobra"
)

func NewRootCmd() *cobra.Command {
	root := &cobra.Command{
		Use:   "envport",
		Short: "Snapshot and restore environment variable sets",
	}

	s, err := store.Default()
	if err != nil {
		panic(err)
	}
	m := store.NewManager(s)

	root.AddCommand(newSaveCmd(m))
	root.AddCommand(newLoadCmd(m))
	root.AddCommand(newListCmd(m))
	root.AddCommand(newDeleteCmd(m))
	root.AddCommand(newRenameCmd(m))
	root.AddCommand(newCopyCmd(m))
	root.AddCommand(newExportCmd(m))
	root.AddCommand(newDiffCmd(m))

	return root
}
