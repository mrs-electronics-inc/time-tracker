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

type projectEditManager interface {
	EditProject(name, newName, code, category string) (*utils.ProjectMutationResult, error)
}

type projectRemoveManager interface {
	RemoveProject(name string) error
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

var projectEditCmd = &cobra.Command{
	Use:   "edit <name>",
	Short: "Edit a project",
	Long:  "Edit a project's name, code, or category.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		newName, err := cmd.Flags().GetString("name")
		if err != nil {
			return fmt.Errorf("failed to parse name flag: %w", err)
		}

		code, err := cmd.Flags().GetString("code")
		if err != nil {
			return fmt.Errorf("failed to parse code flag: %w", err)
		}

		category, err := cmd.Flags().GetString("category")
		if err != nil {
			return fmt.Errorf("failed to parse category flag: %w", err)
		}

		nameChanged := cmd.Flags().Changed("name")
		codeChanged := cmd.Flags().Changed("code")
		categoryChanged := cmd.Flags().Changed("category")

		storage, err := utils.NewFileStorage(config.DataFilePath())
		if err != nil {
			return fmt.Errorf("failed to initialize storage: %w", err)
		}

		taskManager := utils.NewTaskManager(storage)
		return editProject(storage, taskManager, args[0], newName, code, category, nameChanged, codeChanged, categoryChanged, os.Stdout)
	},
}

var projectRemoveCmd = &cobra.Command{
	Use:   "remove <name>",
	Short: "Remove a project",
	Long:  "Remove a project if it is not referenced by any time entries.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		storage, err := utils.NewFileStorage(config.DataFilePath())
		if err != nil {
			return fmt.Errorf("failed to initialize storage: %w", err)
		}

		taskManager := utils.NewTaskManager(storage)
		return removeProject(taskManager, args[0], os.Stdout)
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

func findProjectByName(projects []models.Project, name string) (models.Project, bool) {
	lookupName := strings.TrimSpace(name)
	for _, project := range projects {
		if strings.EqualFold(project.Name, lookupName) {
			return project, true
		}
	}

	return models.Project{}, false
}

func editProject(
	storage projectListStorage,
	taskManager projectEditManager,
	name, newName, code, category string,
	nameChanged, codeChanged, categoryChanged bool,
	out io.Writer,
) error {
	if !nameChanged && !codeChanged && !categoryChanged {
		return fmt.Errorf("at least one flag must be provided: --name, --code, or --category")
	}

	if nameChanged && strings.TrimSpace(newName) == "" {
		return fmt.Errorf("project name cannot be empty")
	}

	if !codeChanged || !categoryChanged {
		projects, err := storage.LoadProjects()
		if err != nil {
			return fmt.Errorf("failed to load projects: %w", err)
		}

		source, found := findProjectByName(projects, name)
		if found {
			if !codeChanged {
				code = source.Code
			}
			if !categoryChanged {
				category = source.Category
			}
		}
	}

	if !nameChanged {
		newName = ""
	}

	result, err := taskManager.EditProject(name, newName, code, category)
	if err != nil {
		return fmt.Errorf("failed to edit project: %w", err)
	}

	if result.Merged {
		fmt.Fprintf(out, "Merged project %q into %q (%d entries rewritten)\n", result.SourceName, result.TargetName, result.RewrittenEntries)
		return nil
	}

	if result.Renamed {
		fmt.Fprintf(out, "Renamed project %q to %q (%d entries rewritten)\n", result.SourceName, result.TargetName, result.RewrittenEntries)
		return nil
	}

	fmt.Fprintf(out, "Updated project %q\n", result.SourceName)
	return nil
}

func removeProject(taskManager projectRemoveManager, name string, out io.Writer) error {
	trimmedName := strings.TrimSpace(name)
	if err := taskManager.RemoveProject(trimmedName); err != nil {
		return fmt.Errorf("failed to remove project: %w", err)
	}

	fmt.Fprintf(out, "Removed project %q\n", trimmedName)
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
	projectEditCmd.Flags().String("name", "", "new project name")
	projectEditCmd.Flags().String("code", "", "external project code")
	projectEditCmd.Flags().String("category", "", "project category")

	projectCmd.AddCommand(projectAddCmd)
	projectCmd.AddCommand(projectEditCmd)
	projectCmd.AddCommand(projectRemoveCmd)
	projectCmd.AddCommand(projectListCmd)
	rootCmd.AddCommand(projectCmd)
}
