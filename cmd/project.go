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

type projectAddManager interface {
	AddProject(name, code, category string) (*models.Project, error)
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

var projectAddCmd = &cobra.Command{
	Use:   "add <name>",
	Short: "Add a project",
	Long:  "Add a project with optional code and category metadata.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		code, err := cmd.Flags().GetString("code")
		if err != nil {
			return fmt.Errorf("failed to parse code flag: %w", err)
		}

		category, err := cmd.Flags().GetString("category")
		if err != nil {
			return fmt.Errorf("failed to parse category flag: %w", err)
		}

		storage, err := utils.NewFileStorage(config.DataFilePath())
		if err != nil {
			return fmt.Errorf("failed to initialize storage: %w", err)
		}

		taskManager := utils.NewTaskManager(storage)
		return addProject(taskManager, args[0], code, category, os.Stdout)
	},
}

func addProject(taskManager projectAddManager, name, code, category string, out io.Writer) error {
	project, err := taskManager.AddProject(name, code, category)
	if err != nil {
		return fmt.Errorf("failed to add project: %w", err)
	}

	fmt.Fprintf(out, "Added project %q\n", project.Name)
	return nil
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
	projectAddCmd.Flags().String("code", "", "external project code")
	projectAddCmd.Flags().String("category", "", "project category")

	projectCmd.AddCommand(projectAddCmd)
	projectCmd.AddCommand(projectListCmd)
	rootCmd.AddCommand(projectCmd)
}
