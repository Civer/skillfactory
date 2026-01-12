// Package system provides system-level commands for HabitWire API
package system

import (
	"encoding/json"
	"fmt"

	"habitwire/client"

	"github.com/spf13/cobra"
)

// HealthResponse represents the health check response
type HealthResponse struct {
	Status  string `json:"status"`
	Version string `json:"version,omitempty"`
}

// ExportData represents exported data
type ExportData struct {
	Habits     json.RawMessage `json:"habits,omitempty"`
	Categories json.RawMessage `json:"categories,omitempty"`
	CheckIns   json.RawMessage `json:"checkins,omitempty"`
}

// RegisterHealthCommand creates the health command
func RegisterHealthCommand(c *client.Client, printJSON func(interface{}) error) *cobra.Command {
	return &cobra.Command{
		Use:   "health",
		Short: "Check API health status",
		RunE: func(cmd *cobra.Command, args []string) error {
			data, err := c.Get("/health")
			if err != nil {
				return err
			}
			var health HealthResponse
			if err := json.Unmarshal(data, &health); err != nil {
				return fmt.Errorf("failed to parse health response: %w", err)
			}
			return printJSON(health)
		},
	}
}

// RegisterExportCommand creates the export command
func RegisterExportCommand(c *client.Client, printJSON func(interface{}) error) *cobra.Command {
	var format string

	cmd := &cobra.Command{
		Use:   "export",
		Short: "Export all data",
		Long:  "Export all habits, categories, and check-ins in JSON or CSV format.",
		RunE: func(cmd *cobra.Command, args []string) error {
			endpoint := "/export"
			if format != "" {
				endpoint += "?format=" + format
			}
			data, err := c.Get(endpoint)
			if err != nil {
				return err
			}
			// For CSV format, just print raw output
			if format == "csv" {
				fmt.Println(string(data))
				return nil
			}
			// For JSON format, parse and re-print
			var export ExportData
			if err := json.Unmarshal(data, &export); err != nil {
				// If parsing fails, just output raw data
				fmt.Println(string(data))
				return nil
			}
			return printJSON(export)
		},
	}
	cmd.Flags().StringVarP(&format, "format", "f", "json", "Export format: json or csv")
	return cmd
}
