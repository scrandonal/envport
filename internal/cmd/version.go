package cmd

import (
	"fmt"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

func newVersionCmd(m Manager) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Manage snapshot version tags",
	}
	cmd.AddCommand(newVersionAddCmd(m))
	cmd.AddCommand(newVersionListCmd(m))
	cmd.AddCommand(newVersionRemoveCmd(m))
	return cmd
}

func newVersionAddCmd(m Manager) *cobra.Command {
	var note string
	cmd := &cobra.Command{
		Use:   "add <snapshot> <tag>",
		Short: "Add a version tag to a snapshot",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := m.AddVersion(args[0], args[1], note); err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "tagged %q as %q\n", args[0], args[1])
			return nil
		},
	}
	cmd.Flags().StringVarP(&note, "note", "n", "", "optional note for the version")
	return cmd
}

func newVersionListCmd(m Manager) *cobra.Command {
	return &cobra.Command{
		Use:   "list <snapshot>",
		Short: "List version tags for a snapshot",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			versions, err := m.ListVersions(args[0])
			if err != nil {
				return err
			}
			if len(versions) == 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "no versions found")
				return nil
			}
			w := tabwriter.NewWriter(cmd.OutOrStdout(), 0, 0, 2, ' ', 0)
			fmt.Fprintln(w, "TAG\tCREATED\tNOTE")
			for _, v := range versions {
				fmt.Fprintf(w, "%s\t%s\t%s\n", v.Tag, v.CreatedAt.Format("2006-01-02 15:04:05"), v.Note)
			}
			return w.Flush()
		},
	}
}

func newVersionRemoveCmd(m Manager) *cobra.Command {
	return &cobra.Command{
		Use:   "remove <snapshot> <tag>",
		Short: "Remove a version tag from a snapshot",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := m.RemoveVersion(args[0], args[1]); err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "removed version tag %q from %q\n", args[1], args[0])
			return nil
		},
	}
}
