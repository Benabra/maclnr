package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

var storageCmd = &cobra.Command{
	Use:   "storage",
	Short: "List all storage devices connected to the system",
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
					err := listStorageDevices(outputFormat)
					if err != nil {
						log.Printf("Failed to list storage devices: %v\n", err)
					}
				}
			}
		} else {
			err := listStorageDevices(outputFormat)
			if err != nil {
				log.Fatalf("Failed to list storage devices: %v\n", err)
			}
		}
	},
}

func listStorageDevices(outputFormat string) error {
	var cmd *exec.Cmd
	if runtime.GOOS == "darwin" {
		cmd = exec.Command("diskutil", "list")
	} else if runtime.GOOS == "linux" {
		cmd = exec.Command("lsblk", "-o", "NAME,FSTYPE,SIZE,MOUNTPOINT")
	} else {
		return fmt.Errorf("unsupported platform")
	}
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("executing %s: %w", cmd.Path, err)
	}

	var devices []map[string]string
	scanner := bufio.NewScanner(strings.NewReader(string(output)))

	if runtime.GOOS == "darwin" {
		var currentDevice map[string]string
		for scanner.Scan() {
			line := scanner.Text()
			if strings.HasPrefix(line, "/dev/") {
				if currentDevice != nil {
					devices = append(devices, currentDevice)
				}
				currentDevice = map[string]string{"Identifier": line}
			} else if strings.Contains(line, "GUID_partition_scheme") || strings.Contains(line, "FDisk_partition_scheme") || strings.Contains(line, "Apple_partition_scheme") {
				currentDevice["Type"] = line
			} else if strings.Contains(line, " (disk") {
				parts := strings.Split(line, " ")
				currentDevice["Name"] = strings.Join(parts[:len(parts)-1], " ")
				currentDevice["Size"] = parts[len(parts)-1]
			}
		}
		if currentDevice != nil {
			devices = append(devices, currentDevice)
		}
	} else {
		// linux lsblk output parsing
		var headers []string
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if line == "" {
				continue
			}
			fields := strings.Fields(line)
			if headers == nil {
				headers = fields
				continue
			}
			device := make(map[string]string)
			for i, h := range headers {
				if i < len(fields) {
					device[h] = fields[i]
				}
			}
			devices = append(devices, device)
		}
	}

	switch outputFormat {
	case "json":
		jsonOutput, _ := json.MarshalIndent(devices, "", "  ")
		fmt.Println(string(jsonOutput))
	case "yaml":
		yamlOutput, _ := yaml.Marshal(devices)
		fmt.Println(string(yamlOutput))
	default:
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Identifier", "Type", "Name", "Size"})
		for _, device := range devices {
			table.Append([]string{device["Identifier"], device["Type"], device["Name"], device["Size"]})
		}
		table.Render()
	}

	return nil
}

func init() {
	scanCmd.AddCommand(storageCmd)
	storageCmd.Flags().StringP("output", "o", "txt", "Output format (txt, json, yaml)")
	storageCmd.Flags().BoolP("watch", "w", false, "Watch the storage devices and refresh every 2 seconds")
}
