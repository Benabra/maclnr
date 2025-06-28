package cmd

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Clean a directory by removing .DS_Store files and files larger than a specified size",
	Run: func(cmd *cobra.Command, args []string) {
		dir, _ := cmd.Flags().GetString("dir")
		dryRun, _ := cmd.Flags().GetBool("dry-run")
		verbose, _ := cmd.Flags().GetBool("verbose")
		removeDS, _ := cmd.Flags().GetBool("ds-store")
		minSize, _ := cmd.Flags().GetInt64("min-size")

		if dir == "" {
			fmt.Println("Please specify a directory with --dir")
			os.Exit(1)
		}

		// Confirm with the user before proceeding
		confirm, _ := cmd.Flags().GetBool("confirm")
		if !confirm {
			reader := bufio.NewReader(os.Stdin)
			fmt.Printf("Are you sure you want to clean the directory %s? (y/N): ", dir)
			response, _ := reader.ReadString('\n')
			response = strings.TrimSpace(response)
			if strings.ToLower(response) != "y" && strings.ToLower(response) != "yes" {
				fmt.Println("Clean operation canceled.")
				os.Exit(0)
			}
		}

		err := cleanDir(dir, dryRun, verbose, removeDS, minSize)
		if err != nil {
			log.Fatalf("Failed to clean directory: %v\n", err)
		}
		fmt.Printf("\033[32mSuccessfully cleaned directory: %s\033[0m\n", dir)
	},
}

func cleanDir(dir string, dryRun bool, verbose bool, removeDS bool, minSize int64) error {
	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("walking %s: %w", path, err)
		}

		// Remove .DS_Store files
		if removeDS && info.Name() == ".DS_Store" {
			if dryRun {
				fmt.Printf("\033[33m[Dry Run]\033[0m Would remove: %s\n", path)
			} else {
				err := os.Remove(path)
				if err != nil {
					return fmt.Errorf("removing %s: %w", path, err)
				}
				if verbose {
					fmt.Printf("\033[32mRemoved: %s\033[0m\n", path)
				}
			}
			return nil
		}

		// Remove files larger than minSize
		if !info.IsDir() && info.Size() >= minSize {
			if dryRun {
				fmt.Printf("\033[33m[Dry Run]\033[0m Would remove: %s (%d bytes)\n", path, info.Size())
			} else {
				err := os.Remove(path)
				if err != nil {
					return fmt.Errorf("removing %s: %w", path, err)
				}
				if verbose {
					fmt.Printf("\033[32mRemoved: %s (%d bytes)\033[0m\n", path, info.Size())
				}
			}
			return nil
		}

		return nil
	})
}

func init() {
	rootCmd.AddCommand(cleanCmd)
	cleanCmd.Flags().StringP("dir", "d", "", "Directory to clean")
	cleanCmd.Flags().Bool("dry-run", false, "Dry run (do not delete files)")
	cleanCmd.Flags().Bool("verbose", false, "Verbose output")
	cleanCmd.Flags().Bool("confirm", false, "Skip confirmation prompt")
	cleanCmd.Flags().Bool("ds-store", false, "Only remove .DS_Store files")
	cleanCmd.Flags().Int64("min-size", 0, "Minimum file size in bytes to clean")
}
