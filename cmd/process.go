package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

var processCmd = &cobra.Command{
	Use:   "process",
	Short: "List all running processes (user, id, memory, CPU usage, name)",
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
					err := listProcesses(outputFormat)
					if err != nil {
						log.Printf("Failed to list processes: %v\n", err)
					}
				}
			}
		} else {
			err := listProcesses(outputFormat)
			if err != nil {
				log.Fatalf("Failed to list processes: %v\n", err)
			}
		}
	},
}

func listProcesses(outputFormat string) error {
	cmd := exec.Command("ps", "aux")
	output, err := cmd.Output()
	if err != nil {
		return err
	}

	lines := strings.Split(string(output), "\n")
	if len(lines) < 1 {
		return fmt.Errorf("no process information available")
	}

	var processes []map[string]string
	for _, line := range lines[1:] {
		if strings.TrimSpace(line) == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) >= 11 {
			processes = append(processes, map[string]string{
				"User":    fields[0],
				"PID":     fields[1],
				"%CPU":    fields[2],
				"%MEM":    fields[3],
				"Command": strings.Join(fields[10:], " "),
			})
		}
	}

	switch outputFormat {
	case "json":
		jsonOutput, _ := json.MarshalIndent(processes, "", "  ")
		fmt.Println(string(jsonOutput))
	case "yaml":
		yamlOutput, _ := yaml.Marshal(processes)
		fmt.Println(string(yamlOutput))
	default:
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"User", "PID", "%CPU", "%MEM", "Command"})
		for _, process := range processes {
			table.Append([]string{process["User"], process["PID"], process["%CPU"], process["%MEM"], process["Command"]})
		}
		table.Render()
	}

	return nil
}

func init() {
	scanCmd.AddCommand(processCmd)
	processCmd.Flags().StringP("output", "o", "txt", "Output format (txt, json, yaml)")
	processCmd.Flags().BoolP("watch", "w", false, "Watch the processes and refresh every 2 seconds")
}
