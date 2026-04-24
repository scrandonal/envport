package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

func newNamespaceCmd(m Manager) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "namespace",
		Short: "Manage snapshot namespaces",
	}
	cmd.AddCommand(
		newNamespaceSetCmd(m),
		newNamespaceGetCmd(m),
		newNamespaceClearCmd(m),
		newNamespaceListCmd(m),
	)
	return cmd
}

func newNamespaceSetCmd(m Manager) *cobra.Command {
	return &cobra.Command{
		Use:   "set <snapshot> <namespace>",
		Short: "Assign a namespace to a snapshot",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := m.SetNamespace(args[0], args[1]); err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "namespace %q set on %q\n", args[1], args[0])
			return nil
		},
	}
}

func newNamespaceGetCmd(m Manager) *cobra.Command {
	return &cobra.Command{
		Use:   "get <snapshot>",
		Short: "Show the namespace of a snapshot",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ns, err := m.GetNamespace(args[0])
			if err != nil {
				return err
			}
			if ns == "" {
				fmt.Fprintln(cmd.OutOrStdout(), "(no namespace set)")
			} else {
				fmt.Fprintln(cmd.OutOrStdout(), ns)
			}
			return nil
		},
	}
}

func newNamespaceClearCmd(m Manager) *cobra.Command {
	return &cobra.Command{
		Use:   "clear <snapshot>",
		Short: "Remove the namespace from a snapshot",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := m.ClearNamespace(args[0]); err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "namespace cleared from %q\n", args[0])
			return nil
		},
	}
}

func newNamespaceListCmd(m Manager) *cobra.Command {
	return &cobra.Command{
		Use:   "list <namespace>",
		Short: "List snapshots in a namespace",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			names, err := m.ListByNamespace(args[0])
			if err != nil {
				return err
			}
			if len(names) == 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "(no snapshots in namespace)")
				return nil
			}
			fmt.Fprintln(cmd.OutOrStdout(), strings.Join(names, "\n"))
			return nil
		},
	}
}
