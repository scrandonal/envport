package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"envport/internal/snapshot"
)

func newImportCmd(m Manager) *cobra.Command {
	var format string
	var overwrite bool

	cmd := &cobra.Command{
		Use:   "import <name>",
		Short: "Import a snapshot from a dotenv or shell export file",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]

			if !overwrite {
				names, err := m.List()
				if err != nil {
					return err
				}
				for _, n := range names {
					if n == name {
						return fmt.Errorf("snapshot %q already exists; use --overwrite to replace", name)
					}
				}
			}

			data, err := os.ReadFile("-")
			if err != nil {
				// not stdin path, read file arg from flag
			}
			_ = data

			input, err := cmd.Flags().GetString("file")
			if err != nil {
				return err
			}

			var raw []byte
			if input == "-" || input == "" {
				raw, err = os.ReadFile("/dev/stdin")
			} else {
				raw, err = os.ReadFile(input)
			}
			if err != nil {
				return fmt.Errorf("reading input: %w", err)
			}

			envMap, err := parseInput(string(raw), format)
			if err != nil {
				return err
			}

			snap := snapshot.New(envMap)
			if err := m.Save(name, snap); err != nil {
				return err
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Imported %d variables into snapshot %q\n", len(envMap), name)
			return nil
		},
	}

	cmd.Flags().StringVarP(&format, "format", "f", "dotenv", "Input format: dotenv or shell")
	cmd.Flags().BoolVar(&overwrite, "overwrite", false, "Overwrite existing snapshot")
	cmd.Flags().String("file", "-", "Input file path (default: stdin)")
	return cmd
}

func parseInput(raw, format string) (map[string]string, error) {
	envMap := make(map[string]string)
	for _, line := range strings.Split(raw, "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		if format == "shell" {
			line = strings.TrimPrefix(line, "export ")
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		val := strings.Trim(strings.TrimSpace(parts[1]), `"`)
		if key != "" {
			envMap[key] = val
		}
	}
	if len(envMap) == 0 {
		return nil, fmt.Errorf("no valid environment variables found in input")
	}
	return envMap, nil
}
