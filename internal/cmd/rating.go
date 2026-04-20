package cmd

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
)

func newRatingCmd(m Manager) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rating",
		Short: "Manage snapshot ratings (1-5)",
	}
	cmd.AddCommand(
		newRatingSetCmd(m),
		newRatingGetCmd(m),
		newRatingClearCmd(m),
	)
	return cmd
}

func newRatingSetCmd(m Manager) *cobra.Command {
	var comment string
	cmd := &cobra.Command{
		Use:   "set <name> <1-5>",
		Short: "Set a rating for a snapshot",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			value, err := strconv.Atoi(args[1])
			if err != nil {
				return fmt.Errorf("invalid rating %q: must be an integer", args[1])
			}
			if err := m.SetRating(args[0], value, comment); err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "rating %d set for %q\n", value, args[0])
			return nil
		},
	}
	cmd.Flags().StringVarP(&comment, "comment", "c", "", "Optional comment")
	return cmd
}

func newRatingGetCmd(m Manager) *cobra.Command {
	return &cobra.Command{
		Use:   "get <name>",
		Short: "Get the rating for a snapshot",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			r, err := m.GetRating(args[0])
			if err != nil {
				return err
			}
			if r == nil {
				fmt.Fprintf(cmd.OutOrStdout(), "no rating set for %q\n", args[0])
				return nil
			}
			if r.Comment != "" {
				fmt.Fprintf(cmd.OutOrStdout(), "%d/5 — %s\n", r.Value, r.Comment)
			} else {
				fmt.Fprintf(cmd.OutOrStdout(), "%d/5\n", r.Value)
			}
			return nil
		},
	}
}

func newRatingClearCmd(m Manager) *cobra.Command {
	return &cobra.Command{
		Use:   "clear <name>",
		Short: "Clear the rating for a snapshot",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := m.ClearRating(args[0]); err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "rating cleared for %q\n", args[0])
			return nil
		},
	}
}
