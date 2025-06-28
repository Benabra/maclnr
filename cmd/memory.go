package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

var memoryCmd = &cobra.Command{
	Use:   "memory",
	Short: "Display the used memory on the system",
	Run: func(cmd *cobra.Command, args []string) {
		outputFormat, _ := cmd.Flags().GetString("output")
		watch, _ := cmd.Flags().GetBool("watch")

		if watch {
			ticker := time.NewTicker(2 * time.Second)
			defer ticker.Stop()
			for {
				select {
				case <-ticker.C:
					clearScreen()
					err := displayMemoryUsage(outputFormat)
					if err != nil {
						log.Printf("Failed to retrieve memory usage: %v\n", err)
					}
				}
			}
		} else {
			err := displayMemoryUsage(outputFormat)
			if err != nil {
				log.Fatalf("Failed to retrieve memory usage: %v\n", err)
			}
		}
	},
}

func displayMemoryUsage(outputFormat string) error {
	var cmd *exec.Cmd

	if runtime.GOOS == "darwin" {
		cmd = exec.Command("vm_stat")
	} else if runtime.GOOS == "linux" {
		cmd = exec.Command("free", "-h")
	} else {
		return fmt.Errorf("unsupported platform")
	}

	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("executing %s: %w", cmd.Path, err)
	}

	if runtime.GOOS == "darwin" {
		return parseAndDisplayMacMemory(string(output), outputFormat)
	} else if runtime.GOOS == "linux" {
		switch outputFormat {
		case "json":
			data := parseLinuxMemoryOutput(string(output))
			jsonOutput, _ := json.MarshalIndent(data, "", "  ")
			fmt.Println(string(jsonOutput))
		case "yaml":
			data := parseLinuxMemoryOutput(string(output))
			yamlOutput, _ := yaml.Marshal(data)
			fmt.Println(string(yamlOutput))
		default:
			fmt.Println(strings.TrimSpace(string(output)))
		}
	}

	return nil
}

func parseLinuxMemoryOutput(output string) map[string]interface{} {
	lines := strings.Split(output, "\n")
	if len(lines) < 2 {
		return nil
	}

	headers := strings.Fields(lines[0])
	values := strings.Fields(lines[1])

	data := make(map[string]interface{})
	for i, header := range headers {
		if i < len(values) {
			data[header] = values[i]
		}
	}

	return data
}

func parseAndDisplayMacMemory(output string, outputFormat string) error {
	scanner := bufio.NewScanner(strings.NewReader(output))
	var pageSize int64
	var memStats = make(map[string]int64)

	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "Mach Virtual Memory Statistics") {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}

		if fields[0] == "Pagesize:" {
			size, err := strconv.ParseInt(fields[1], 10, 64)
			if err == nil {
				pageSize = size
			}
			continue
		}

		key := strings.TrimSuffix(fields[0], ":")
		valueStr := strings.TrimSuffix(fields[1], ".")
		value, err := strconv.ParseInt(valueStr, 10, 64)
		if err == nil {
			memStats[key] = value
		}
	}

	if pageSize == 0 {
		pageSize = 4096 // default page size
	}

	result := make(map[string]interface{})
	for key, value := range memStats {
		size := value * pageSize
		result[key] = fmt.Sprintf("%d bytes", size)
	}

	switch outputFormat {
	case "json":
		jsonOutput, _ := json.MarshalIndent(result, "", "  ")
		fmt.Println(string(jsonOutput))
	case "yaml":
		yamlOutput, _ := yaml.Marshal(result)
		fmt.Println(string(yamlOutput))
	default:
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Type", "Size"})
		for key, value := range result {
			table.Append([]string{key, value.(string)})
		}
		table.Render()
	}

	return nil
}

func init() {
	scanCmd.AddCommand(memoryCmd)
	memoryCmd.Flags().StringP("output", "o", "txt", "Output format (txt, json, yaml)")
	memoryCmd.Flags().BoolP("watch", "w", false, "Watch the memory usage and refresh every 2 seconds")
}
