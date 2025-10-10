package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/opd-ai/asset-generator/pkg/client"
	"github.com/opd-ai/asset-generator/pkg/output"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// statusCmd represents the status command
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Check SwarmUI server status",
	Long: `Query the SwarmUI server and display its current status, including:
  - Server connectivity and response time
  - Available backends and their states
  - Current session information
  - Model loading status

Examples:
  # Check server status
  asset-generator status
  
  # Check status with JSON output
  asset-generator status --format json
  
  # Check status with verbose details
  asset-generator status -v`,
	RunE: runStatus,
}

func init() {
	rootCmd.AddCommand(statusCmd)
}

func runStatus(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	// Get server status
	status, err := assetClient.GetServerStatus(ctx)
	if err != nil {
		return fmt.Errorf("failed to get server status: %w", err)
	}

	// Format output based on user preference
	format := viper.GetString("format")
	outputPath := viper.GetString("output")

	formatter := output.NewFormatter(format)
	var result string

	switch format {
	case "json", "yaml":
		result, err = formatter.Format(status)
		if err != nil {
			return fmt.Errorf("failed to format output: %w", err)
		}
	default:
		// Table/human-readable format
		result = formatStatusTable(status)
	}

	// Write output
	if outputPath != "" {
		if err := os.WriteFile(outputPath, []byte(result), 0644); err != nil {
			return fmt.Errorf("failed to write output file: %w", err)
		}
		if !quiet {
			fmt.Printf("Status written to %s\n", outputPath)
		}
	} else {
		fmt.Print(result)
	}

	return nil
}

func formatStatusTable(status *client.ServerStatus) string {
	var result string

	// Server info
	result += fmt.Sprintf("SwarmUI Server Status\n")
	result += fmt.Sprintf("═══════════════════════════════════════════════\n\n")
	result += fmt.Sprintf("Server URL:      %s\n", status.ServerURL)
	result += fmt.Sprintf("Status:          %s\n", colorizeStatus(status.Status))
	result += fmt.Sprintf("Response Time:   %s\n", status.ResponseTime)
	result += fmt.Sprintf("Version:         %s\n", valueOrNA(status.Version))

	// Session info
	if status.SessionID != "" {
		result += fmt.Sprintf("\nSession\n")
		result += fmt.Sprintf("───────────────────────────────────────────────\n")
		result += fmt.Sprintf("Session ID:      %s\n", status.SessionID)
	}

	// Backends info
	if len(status.Backends) > 0 {
		result += fmt.Sprintf("\nBackends\n")
		result += fmt.Sprintf("───────────────────────────────────────────────\n")
		for _, backend := range status.Backends {
			result += fmt.Sprintf("  • %s\n", backend.ID)
			result += fmt.Sprintf("    Type:          %s\n", backend.Type)
			result += fmt.Sprintf("    Status:        %s\n", colorizeStatus(backend.Status))
			if backend.ModelLoaded != "" {
				result += fmt.Sprintf("    Model Loaded:  %s\n", backend.ModelLoaded)
			}
			if backend.GPU != "" {
				result += fmt.Sprintf("    GPU:           %s\n", backend.GPU)
			}
			result += fmt.Sprintf("\n")
		}
	}

	// Model info
	if status.ModelsCount > 0 {
		result += fmt.Sprintf("Models\n")
		result += fmt.Sprintf("───────────────────────────────────────────────\n")
		result += fmt.Sprintf("Total Available: %d\n", status.ModelsCount)
		if status.ModelsLoaded > 0 {
			result += fmt.Sprintf("Currently Loaded: %d\n", status.ModelsLoaded)
		}
	}

	// System info
	if status.SystemInfo != nil && len(status.SystemInfo) > 0 {
		result += fmt.Sprintf("\nSystem Information\n")
		result += fmt.Sprintf("───────────────────────────────────────────────\n")
		for key, value := range status.SystemInfo {
			result += fmt.Sprintf("%-20s %v\n", key+":", value)
		}
	}

	return result
}

func colorizeStatus(status string) string {
	// Simple status coloring using ANSI codes (works on Linux terminals)
	switch status {
	case "running", "loaded", "active", "online", "ready":
		return fmt.Sprintf("\033[32m%s\033[0m", status) // Green
	case "idle", "unloaded":
		return fmt.Sprintf("\033[33m%s\033[0m", status) // Yellow
	case "error", "failed", "offline":
		return fmt.Sprintf("\033[31m%s\033[0m", status) // Red
	default:
		return status
	}
}

func valueOrNA(value string) string {
	if value == "" {
		return "N/A"
	}
	return value
}
