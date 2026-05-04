// main is the entry point for the envport CLI tool.
// It wires together the root command and executes it.
package main

import (
	"fmt"
	"os"

	"github.com/user/envport/internal/cmd"
)

func main() {
	root := cmd.NewRootCmd()
	if err := root.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
