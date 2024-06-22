package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

// FileInfo structure to store file path and size
type FileInfo struct {
	Path string
	Size int64
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all files in a directory recursively, ordered by size from big to small",
	Run: func(cmd *cobra.Command, args []string) {
		dir, _ := cmd.Flags().GetString("dir")
		minSize, _ := cmd.Flags().GetInt64("min-size")
		outputFormat, _ := cmd.Flags().GetString("output")
		watch, _ := cmd.Flags().GetBool("watch")

		if dir == "" {
			fmt.Println("Please specify a directory with --dir")
			os.Exit(1)
		}

		if watch {
			ticker := time.NewTicker(2 * time.Second)
			defer ticker.Stop()
			for {
				select {
				case <-ticker.C:
					clearScreen()
					err := listFiles(dir, minSize, outputFormat)
					if err != nil {
						log.Printf("Failed to list files in directory: %v\n", err)
					}
				}
			}
		} else {
			err := listFiles(dir, minSize, outputFormat)
			if err != nil {
				log.Fatalf("Failed to list files in directory: %v\n", err)
			}
		}
	},
}

func listFiles(dir string, minSize int64, outputFormat string) error {
	files, err := listFilesBySize(dir, minSize)
	if err != nil {
		return err
	}

	switch outputFormat {
	case "json":
		jsonOutput, _ := json.MarshalIndent(files, "", "  ")
		fmt.Println(string(jsonOutput))
	case "yaml":
		yamlOutput, _ := yaml.Marshal(files)
		fmt.Println(string(yamlOutput))
	default:
		displayFilesTable(files)
	}

	return nil
}

func listFilesBySize(dir string, minSize int64) ([]FileInfo, error) {
	var files []FileInfo

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && info.Size() >= minSize {
			files = append(files, FileInfo{Path: path, Size: info.Size()})
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	// Sort files by size in descending order
	sort.Slice(files, func(i, j int) bool {
		return files[i].Size > files[j].Size
	})

	return files, nil
}

func displayFilesTable(files []FileInfo) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Path", "Size (bytes)"})

	for _, file := range files {
		table.Append([]string{file.Path, fmt.Sprintf("%d", file.Size)})
	}

	table.Render()
}

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.Flags().StringP("dir", "d", "", "Directory to list files")
	listCmd.Flags().Int64("min-size", 0, "Minimum file size in bytes")
	listCmd.Flags().StringP("output", "o", "txt", "Output format (txt, json, yaml)")
	listCmd.Flags().BoolP("watch", "w", false, "Watch the directory and refresh every 2 seconds")
}
