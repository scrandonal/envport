package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newBookmarkCmd(m Manager) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "bookmark",
		Short: "Manage snapshot bookmarks",
	}
	cmd.AddCommand(newBookmarkAddCmd(m))
	cmd.AddCommand(newBookmarkRemoveCmd(m))
	cmd.AddCommand(newBookmarkListCmd(m))
	return cmd
}

func newBookmarkAddCmd(m Manager) *cobra.Command {
	var label string
	cmd := &cobra.Command{
		Use:   "add <name>",
		Short: "Bookmark a snapshot",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := m.AddBookmark(args[0], label); err != nil {
				return fmt.Errorf("bookmark add: %w", err)
			}
			fmt.Fprintf(cmd.OutOrStdout(), "bookmarked %q\n", args[0])
			return nil
		},
	}
	cmd.Flags().StringVarP(&label, "label", "l", "", "optional label")
	return cmd
}

func newBookmarkRemoveCmd(m Manager) *cobra.Command {
	return &cobra.Command{
		Use:   "remove <name>",
		Short: "Remove a bookmark",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := m.RemoveBookmark(args[0]); err != nil {
				return fmt.Errorf("bookmark remove: %w", err)
			}
			fmt.Fprintf(cmd.OutOrStdout(), "removed bookmark %q\n", args[0])
			return nil
		},
	}
}

func newBookmarkListCmd(m Manager) *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List bookmarks",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			bm, err := m.ListBookmarks()
			if err != nil {
				return err
			}
			if len(bm) == 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "no bookmarks")
				return nil
			}
			for _, b := range bm {
				if b.Label != "" {
					fmt.Fprintf(cmd.OutOrStdout(), "%s\t%s\n", b.Name, b.Label)
				} else {
					fmt.Fprintln(cmd.OutOrStdout(), b.Name)
				}
			}
			return nil
		},
	}
}
