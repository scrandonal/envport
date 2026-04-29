package cmd

import (
	"fmt"
	"os"
	"runtime"
	"text/tabwriter"

	"github.com/spf13/cobra"

	store "envport/internal/store"
)

func newRuntimeCmd(m Manager) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "runtime",
		Short: "Manage runtime metadata for snapshots",
	}
	cmd.AddCommand(newRuntimeCaptureCmd(m))
	cmd.AddCommand(newRuntimeShowCmd(m))
	cmd.AddCommand(newRuntimeClearCmd(m))
	return cmd
}

func newRuntimeCaptureCmd(m Manager) *cobra.Command {
	return &cobra.Command{
		Use:   "capture <name>",
		Short: "Capture current runtime info for a snapshot",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			host, _ := os.Hostname()
			user := os.Getenv("USER")
			shell := os.Getenv("SHELL")
			info := store.RuntimeInfo{
				OS:    runtime.GOOS,
				Arch:  runtime.GOARCH,
				Host:  host,
				User:  user,
				Shell: shell,
			}
			if err := m.SetRuntime(args[0], info); err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "runtime captured for %q\n", args[0])
			return nil
		},
	}
}

func newRuntimeShowCmd(m Manager) *cobra.Command {
	return &cobra.Command{
		Use:   "show <name>",
		Short: "Show runtime info for a snapshot",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			info, err := m.GetRuntime(args[0])
			if err != nil {
				return err
			}
			w := tabwriter.NewWriter(cmd.OutOrStdout(), 0, 0, 2, ' ', 0)
			fmt.Fprintf(w, "OS:\t%s\n", info.OS)
			fmt.Fprintf(w, "Arch:\t%s\n", info.Arch)
			fmt.Fprintf(w, "Host:\t%s\n", info.Host)
			fmt.Fprintf(w, "User:\t%s\n", info.User)
			fmt.Fprintf(w, "Shell:\t%s\n", info.Shell)
			return w.Flush()
		},
	}
}

func newRuntimeClearCmd(m Manager) *cobra.Command {
	return &cobra.Command{
		Use:   "clear <name>",
		Short: "Clear runtime info for a snapshot",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := m.ClearRuntime(args[0]); err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "runtime cleared for %q\n", args[0])
			return nil
		},
	}
}
