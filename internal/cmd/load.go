package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"envport/internal/store"
)

func newLoadCmd(mgr *store.Manager) *cobra.Command {
	var shellFmt string
	cmd := &cobra.Command{
		Use:   "load <name>",
		Short: "Print shell export statements for a saved snapshot",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]
			snap, err := mgr.Load(name)
			if err != nil {
				return fmt.Errorf("load: %w", err)
			}
			if err := validateShell(shellFmt); err != nil {
				return err
			}
			out := cmd.OutOrStdout()
			for k, v := range snap.Vars {
				if err := printExport(out, shellFmt, k, v); err != nil {
					return fmt.Errorf("load: writing output: %w", err)
				}
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&shellFmt, "shell", "bash", "Shell format: bash, fish")
	return cmd
}

// validateShell returns an error if the given shell format is not supported.
func validateShell(shell string) error {
	switch shell {
	case "bash", "fish":
		return nil
	default:
		return fmt.Errorf("unsupported shell format %q: must be one of: bash, fish", shell)
	}
}

func printExport(out interface{ Write([]byte) (int, error) }, shell, k, v string) error {
	var err error
	switch shell {
	case "fish":
		_, err = fmt.Fprintf(out, "set -x %s %q\n", k, v)
	default:
		_, err = fmt.Fprintf(out, "export %s=%q\n", k, v)
	}
	return err
}
