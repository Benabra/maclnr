package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "maclnr",
	Short: "A command-line tool to manage directories",
	Long: `maclnr provides utilities to list large files, clean directories,
and inspect system information. Most commands support a --watch mode
for continuous updates.`,
	Version: version,
}

var version = "v0.1.0"

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(scanCmd)
}
