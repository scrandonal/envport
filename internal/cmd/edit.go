package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

func newEditCmd(m Manager) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "edit <name>",
		Short: "Open a snapshot in your default editor",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]

			snap, err := m.Load(name)
			if err != nil {
				return fmt.Errorf("snapshot %q not found", name)
			}

			tmp, err := os.CreateTemp("", "envport-edit-*.env")
			if err != nil {
				return fmt.Errorf("failed to create temp file: %w", err)
			}
			defer os.Remove(tmp.Name())

			for k, v := range snap.Vars {
				fmt.Fprintf(tmp, "%s=%s\n", k, v)
			}
			tmp.Close()

			editor := os.Getenv("EDITOR")
			if editor == "" {
				editor = "vi"
			}

			c := exec.Command(editor, tmp.Name())
			c.Stdin = os.Stdin
			c.Stdout = os.Stdout
			c.Stderr = os.Stderr
			if err := c.Run(); err != nil {
				return fmt.Errorf("editor exited with error: %w", err)
			}

			updated, err := parseInput(tmp.Name(), "dotenv")
			if err != nil {
				return fmt.Errorf("failed to parse edited file: %w", err)
			}

			if err := m.Save(name, updated, true); err != nil {
				return fmt.Errorf("failed to save snapshot: %w", err)
			}

			fmt.Fprintf(cmd.OutOrStdout(), "snapshot %q updated\n", name)
			return nil
		},
	}
	return cmd
}
