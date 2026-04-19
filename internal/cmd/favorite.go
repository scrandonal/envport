package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newFavoriteCmd(m Manager) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "favorite",
		Short: "Manage favorite snapshots",
	}
	cmd.AddCommand(newFavoriteAddCmd(m))
	cmd.AddCommand(newFavoriteRemoveCmd(m))
	cmd.AddCommand(newFavoriteListCmd(m))
	return cmd
}

func newFavoriteAddCmd(m Manager) *cobra.Command {
	return &cobra.Command{
		Use:   "add <name>",
		Short: "Mark a snapshot as favorite",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := m.AddFavorite(args[0]); err != nil {
				return fmt.Errorf("add favorite: %w", err)
			}
			fmt.Fprintf(cmd.OutOrStdout(), "marked %q as favorite\n", args[0])
			return nil
		},
	}
}

func newFavoriteRemoveCmd(m Manager) *cobra.Command {
	return &cobra.Command{
		Use:   "remove <name>",
		Short: "Unmark a snapshot as favorite",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := m.RemoveFavorite(args[0]); err != nil {
				return fmt.Errorf("remove favorite: %w", err)
			}
			fmt.Fprintf(cmd.OutOrStdout(), "removed %q from favorites\n", args[0])
			return nil
		},
	}
}

func newFavoriteListCmd(m Manager) *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List favorite snapshots",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			favs, err := m.ListFavorites()
			if err != nil {
				return err
			}
			if len(favs) == 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "no favorites")
				return nil
			}
			for _, f := range favs {
				fmt.Fprintln(cmd.OutOrStdout(), f)
			}
			return nil
		},
	}
}
