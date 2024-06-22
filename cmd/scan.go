package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var scanCmd = &cobra.Command{
	Use:   "scan [memory|storage|process]",
	Short: "Scan and display system information",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		switch args[0] {
		case "memory":
			memoryCmd.Run(cmd, args)
		case "storage":
			storageCmd.Run(cmd, args)
		case "process":
			processCmd.Run(cmd, args)
		default:
			fmt.Println("Invalid argument. Use 'memory', 'storage', or 'process'.")
		}
	},
}

func init() {
	rootCmd.AddCommand(scanCmd)
}
