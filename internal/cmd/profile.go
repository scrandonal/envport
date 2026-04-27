package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"envport/internal/store"
)

func newProfileCmd(mgr Manager) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "profile",
		Short: "Manage snapshot profiles (name, description, author)",
	}
	cmd.AddCommand(
		newProfileSetCmd(mgr),
		newProfileGetCmd(mgr),
		newProfileClearCmd(mgr),
	)
	return cmd
}

func newProfileSetCmd(mgr Manager) *cobra.Command {
	var description, author string
	cmd := &cobra.Command{
		Use:   "set <snapshot>",
		Short: "Set profile metadata for a snapshot",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			p := store.Profile{
				Name:        args[0],
				Description: description,
				Author:      author,
			}
			if err := mgr.SetProfile(args[0], p); err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "profile set for %q\n", args[0])
			return nil
		},
	}
	cmd.Flags().StringVar(&description, "description", "", "Profile description")
	cmd.Flags().StringVar(&author, "author", "", "Profile author")
	return cmd
}

func newProfileGetCmd(mgr Manager) *cobra.Command {
	return &cobra.Command{
		Use:   "get <snapshot>",
		Short: "Get profile metadata for a snapshot",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			p, err := mgr.GetProfile(args[0])
			if err != nil {
				return err
			}
			if p.Name == "" && p.Description == "" && p.Author == "" {
				fmt.Fprintf(cmd.OutOrStdout(), "no profile set for %q\n", args[0])
				return nil
			}
			fmt.Fprintf(cmd.OutOrStdout(), "name:        %s\n", p.Name)
			fmt.Fprintf(cmd.OutOrStdout(), "description: %s\n", p.Description)
			fmt.Fprintf(cmd.OutOrStdout(), "author:      %s\n", p.Author)
			return nil
		},
	}
}

func newProfileClearCmd(mgr Manager) *cobra.Command {
	return &cobra.Command{
		Use:   "clear <snapshot>",
		Short: "Clear profile metadata for a snapshot",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := mgr.ClearProfile(args[0]); err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "profile cleared for %q\n", args[0])
			return nil
		},
	}
}
