package cmd

import (
	"fmt"
	"io"
	"os"
	"sort"
	"strings"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"time-tracker/config"
	"time-tracker/models"
	"time-tracker/utils"
)

type projectListStorage interface {
	LoadProjects() ([]models.Project, error)
}

var projectCmd = &cobra.Command{
	Use:   "project",
	Short: "Manage projects",
}

var projectListCmd = &cobra.Command{
	Use:   "list",
	Short: "List projects",
	Long:  "List all projects with metadata.",
	RunE: func(cmd *cobra.Command, args []string) error {
		storage, err := utils.NewFileStorage(config.DataFilePath())
		if err != nil {
			return fmt.Errorf("failed to initialize storage: %w", err)
		}

		return listProjects(storage, os.Stdout)
	},
}

func listProjects(storage projectListStorage, out io.Writer) error {
	projects, err := storage.LoadProjects()
	if err != nil {
		return fmt.Errorf("failed to load projects: %w", err)
	}

	sort.Slice(projects, func(i, j int) bool {
		leftName := strings.ToLower(projects[i].Name)
		rightName := strings.ToLower(projects[j].Name)

		if leftName == rightName {
			return projects[i].Name < projects[j].Name
		}

		return leftName < rightName
	})

	table := tablewriter.NewWriter(out)
	table.SetHeader([]string{"Name", "Code", "Category"})
	table.SetAutoFormatHeaders(false)
	table.SetBorder(true)
	table.SetRowLine(true)
	table.SetAutoWrapText(false)

	for _, project := range projects {
		table.Append([]string{project.Name, project.Code, project.Category})
	}

	table.Render()
	return nil
}

func init() {
	projectCmd.AddCommand(projectListCmd)
	rootCmd.AddCommand(projectCmd)
}
