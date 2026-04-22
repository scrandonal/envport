package cmd

import (
	"fmt"
	"os"
	"runtime"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

func newEnvironmentCmd(m Manager) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "environment",
		Short: "Manage capture-time environment context for snapshots",
	}
	cmd.AddCommand(newEnvironmentCaptureCmd(m))
	cmd.AddCommand(newEnvironmentShowCmd(m))
	cmd.AddCommand(newEnvironmentClearCmd(m))
	return cmd
}

func newEnvironmentCaptureCmd(m Manager) *cobra.Command {
	return &cobra.Command{
		Use:   "capture <name>",
		Short: "Capture the current system environment context for a snapshot",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			hostname, _ := os.Hostname()
			user := os.Getenv("USER")
			if user == "" {
				user = os.Getenv("USERNAME")
			}
			shell := os.Getenv("SHELL")
			if err := m.SetEnvironment(args[0], environmentRecord(hostname, user, runtime.GOOS, shell)); err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "environment context captured for %q\n", args[0])
			return nil
		},
	}
}

func newEnvironmentShowCmd(m Manager) *cobra.Command {
	return &cobra.Command{
		Use:   "show <name>",
		Short: "Show the captured environment context for a snapshot",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			rec, err := m.GetEnvironment(args[0])
			if err != nil {
				return err
			}
			w := tabwriter.NewWriter(cmd.OutOrStdout(), 0, 0, 2, ' ', 0)
			fmt.Fprintf(w, "hostname:\t%s\n", rec.Hostname)
			fmt.Fprintf(w, "user:\t%s\n", rec.User)
			fmt.Fprintf(w, "os:\t%s\n", rec.OS)
			fmt.Fprintf(w, "shell:\t%s\n", rec.Shell)
			return w.Flush()
		},
	}
}

func newEnvironmentClearCmd(m Manager) *cobra.Command {
	return &cobra.Command{
		Use:   "clear <name>",
		Short: "Remove the captured environment context for a snapshot",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := m.ClearEnvironment(args[0]); err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "environment context cleared for %q\n", args[0])
			return nil
		},
	}
}

type envRec = store.EnvironmentRecord

func environmentRecord(hostname, user, goos, shell string) store.EnvironmentRecord {
	return store.EnvironmentRecord{
		Hostname: hostname,
		User:     user,
		OS:       goos,
		Shell:    shell,
	}
}
