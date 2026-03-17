package cmd

import (
	"fmt"
	"os"
	"strings"
	"time"
	"time-tracker/config"
	"time-tracker/models"
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

		category, err := cmd.Flags().GetString("category")
		if err != nil {
			return fmt.Errorf("failed to parse category flag: %w", err)
		}
		categoryProvided := cmd.Flags().Changed("category")

		days, err := cmd.Flags().GetInt("days")
		if err != nil {
			return fmt.Errorf("failed to parse days flag: %w", err)
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

		exportData, err := buildExportData(storage, format, category, categoryProvided, days, time.Now())
		if err != nil {
			return err
		}

		// Write to output destination
		if output == "" {
			// Write to stdout without adding extra newline to preserve exact TSV format
			_, err := os.Stdout.WriteString(exportData)
			if err != nil {
				return fmt.Errorf("failed to write to stdout: %w", err)
			}
		} else {
			// Write to file with restricted permissions (owner read/write only)
			// for privacy of time tracking data
			err := os.WriteFile(output, []byte(exportData), 0600)
			if err != nil {
				return fmt.Errorf("failed to write to file %q: %w", output, err)
			}
			fmt.Fprintf(os.Stderr, "Data exported to %s\n", output)
		}

		return nil
	},
}

func init() {
	exportCmd.Flags().StringP("format", "f", "daily-projects", "Export format: \"daily-projects\" or \"raw\"")
	exportCmd.Flags().StringP("output", "o", "", "Output file path (default: stdout)")
	exportCmd.Flags().String("category", "", "Filter exported rows by project category (case-insensitive)")
	exportCmd.Flags().IntP("days", "d", 7, "Number of past days to include in export")

	rootCmd.AddCommand(exportCmd)
}

type exportStorage interface {
	Load() ([]models.TimeEntry, error)
	LoadProjects() ([]models.Project, error)
}

func buildExportData(storage exportStorage, format, category string, categoryProvided bool, days int, now time.Time) (string, error) {
	entries, err := storage.Load()
	if err != nil {
		return "", fmt.Errorf("failed to load entries: %w", err)
	}

	if days <= 0 {
		return "", fmt.Errorf("days must be a positive integer")
	}
	entries = filterEntriesByPastDays(entries, days, now)

	trimmedCategory := strings.TrimSpace(category)
	if categoryProvided && trimmedCategory == "" {
		return "", fmt.Errorf("category cannot be empty or whitespace")
	}

	switch format {
	case "daily-projects":
		aggregated := utils.AggregateByProjectDate(entries)
		projects, err := storage.LoadProjects()
		if err != nil {
			return "", fmt.Errorf("failed to load projects: %w", err)
		}
		aggregated = utils.ApplyProjectMetadata(aggregated, projects)
		if categoryProvided {
			aggregated = filterAggregatedByCategory(aggregated, trimmedCategory)
		}

		exportData, err := utils.ExportDailyProjects(aggregated)
		if err != nil {
			return "", fmt.Errorf("failed to export daily projects: %w", err)
		}
		return exportData, nil

	case "raw":
		if categoryProvided {
			projects, err := storage.LoadProjects()
			if err != nil {
				return "", fmt.Errorf("failed to load projects: %w", err)
			}
			entries = filterEntriesByCategory(entries, projects, trimmedCategory)
		}

		exportData, err := utils.ExportRaw(entries)
		if err != nil {
			return "", fmt.Errorf("failed to export raw data: %w", err)
		}
		return exportData, nil

	default:
		return "", fmt.Errorf("invalid format %q. Must be 'daily-projects' or 'raw'", format)
	}
}

func filterEntriesByPastDays(entries []models.TimeEntry, days int, now time.Time) []models.TimeEntry {
	filtered := make([]models.TimeEntry, 0, len(entries))
	for _, entry := range entries {
		daysDiff := int(now.Sub(entry.Start).Hours() / 24)
		if daysDiff >= 0 && daysDiff < days {
			filtered = append(filtered, entry)
		}
	}

	return filtered
}

func filterAggregatedByCategory(entries []utils.ProjectDateEntry, category string) []utils.ProjectDateEntry {
	filtered := make([]utils.ProjectDateEntry, 0, len(entries))
	for _, entry := range entries {
		if strings.EqualFold(entry.ProjectCategory, category) {
			filtered = append(filtered, entry)
		}
	}

	return filtered
}

func filterEntriesByCategory(entries []models.TimeEntry, projects []models.Project, category string) []models.TimeEntry {
	byName := make(map[string]models.Project, len(projects))
	for _, project := range projects {
		byName[project.Name] = project
	}

	filtered := make([]models.TimeEntry, 0, len(entries))
	for _, entry := range entries {
		project, ok := byName[entry.Project]
		if !ok {
			continue
		}
		if strings.EqualFold(project.Category, category) {
			filtered = append(filtered, entry)
		}
	}

	return filtered
}
