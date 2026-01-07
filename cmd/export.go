package cmd

import (
	"fmt"
	"os"
	"time-tracker/config"
	"time-tracker/utils"

	"github.com/spf13/cobra"
)

var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "Export time tracking data to TSV format",
	Long: `Export time tracking data as tab-separated values (TSV) for import into spreadsheet software or other tools.

Supports two export formats:
- daily-projects (default): Aggregated by project and date with combined task descriptions
- raw: Individual time entries with start/end times

Running entries (without end times) and blank entries are excluded from exports.
By default, output is written to stdout. Use --output to write to a file.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		format, err := cmd.Flags().GetString("format")
		if err != nil {
			return fmt.Errorf("failed to parse format flag: %w", err)
		}

		output, err := cmd.Flags().GetString("output")
		if err != nil {
			return fmt.Errorf("failed to parse output flag: %w", err)
		}

		// Validate format
		if format != "daily-projects" && format != "raw" {
			return fmt.Errorf("invalid format %q. Must be 'daily-projects' or 'raw'", format)
		}

		// Load data
		storage, err := utils.NewFileStorage(config.DataFilePath())
		if err != nil {
			return fmt.Errorf("failed to initialize storage: %w", err)
		}

		entries, err := storage.Load()
		if err != nil {
			return fmt.Errorf("failed to load entries: %w", err)
		}

		var exportData string

		switch format {
		case "daily-projects":
			aggregated := utils.AggregateByProjectDate(entries)
			exportData = utils.ExportDailyProjects(aggregated)

		case "raw":
			exportData = utils.ExportRaw(entries)
		}

		// Write to output destination
		if output == "" {
			// Write to stdout
			fmt.Print(exportData)
		} else {
			// Write to file
			err := os.WriteFile(output, []byte(exportData), 0644)
			if err != nil {
				return fmt.Errorf("failed to write to file %q: %w", output, err)
			}
			fmt.Fprintf(os.Stderr, "Data exported to %s\n", output)
		}

		return nil
	},
}

func init() {
	exportCmd.Flags().StringP("format", "f", "daily-projects", "Export format: 'daily-projects' (default) or 'raw'")
	exportCmd.Flags().StringP("output", "o", "", "Output file path (default: stdout)")

	rootCmd.AddCommand(exportCmd)
}
