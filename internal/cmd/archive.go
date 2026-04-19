package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"envport/internal/store"
)

func newArchiveCmd(m Manager) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "archive",
		Short: "Manage snapshot archives",
	}
	cmd.AddCommand(newArchiveSaveCmd(m))
	cmd.AddCommand(newArchiveListCmd(m))
	cmd.AddCommand(newArchiveClearCmd(m))
	return cmd
}

func newArchiveSaveCmd(m Manager) *cobra.Command {
	return &cobra.Command{
		Use:   "save <name>",
		Short: "Archive a snapshot by name",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]
			snap, err := m.Load(name)
			if err != nil {
				return fmt.Errorf("snapshot %q not found", name)
			}
			s, err := store.Default()
			if err != nil {
				return err
			}
			if err := store.ArchiveSnapshot(s.Base, name, snap.Vars); err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "archived %q\n", name)
			return nil
		},
	}
}

func newArchiveListCmd(m Manager) *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List archived snapshots",
		RunE: func(cmd *cobra.Command, args []string) error {
			s, err := store.Default()
			if err != nil {
				return err
			}
			entries, err := store.ListArchive(s.Base)
			if err != nil {
				return err
			}
			if len(entries) == 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "no archived snapshots")
				return nil
			}
			w := tabwriter.NewWriter(cmd.OutOrStdout(), 0, 0, 2, ' ', 0)
			fmt.Fprintln(w, "NAME\tARCHIVED AT")
			for _, e := range entries {
				fmt.Fprintf(w, "%s\t%s\n", e.Name, e.ArchivedAt.Format("2006-01-02 15:04:05"))
			}
			return w.Flush()
		},
	}
}

func newArchiveClearCmd(_ Manager) *cobra.Command {
	return &cobra.Command{
		Use:   "clear",
		Short: "Clear all archived snapshots",
		RunE: func(cmd *cobra.Command, args []string) error {
			s, err := store.Default()
			if err != nil {
				return err
			}
			if err := store.ClearArchive(s.Base); err != nil {
				return err
			}
			fmt.Fprintln(cmd.OutOrStdout(), "archive cleared")
			_ = os.Stderr
			return nil
		},
	}
}
