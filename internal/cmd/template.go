package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

func newTemplateCmd(m Manager) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "template",
		Short: "Manage snapshot templates",
	}
	cmd.AddCommand(newTemplateListCmd(m))
	cmd.AddCommand(newTemplateApplyCmd(m))
	cmd.AddCommand(newTemplateDeleteCmd(m))
	return cmd
}

func newTemplateListCmd(m Manager) *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List saved templates",
		RunE: func(cmd *cobra.Command, args []string) error {
			names, err := m.ListTemplates()
			if err != nil {
				return err
			}
			if len(names) == 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "no templates saved")
				return nil
			}
			fmt.Fprintln(cmd.OutOrStdout(), strings.Join(names, "\n"))
			return nil
		},
	}
}

func newTemplateApplyCmd(m Manager) *cobra.Command {
	return &cobra.Command{
		Use:   "apply <template> <snapshot>",
		Short: "Create a snapshot from a template",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := m.ApplyTemplate(args[0], args[1]); err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "created snapshot %q from template %q\n", args[1], args[0])
			return nil
		},
	}
}

func newTemplateDeleteCmd(m Manager) *cobra.Command {
	return &cobra.Command{
		Use:   "delete <template>",
		Short: "Delete a template",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := m.DeleteTemplate(args[0]); err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "deleted template %q\n", args[0])
			return nil
		},
	}
}
