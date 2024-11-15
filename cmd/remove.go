/*
Copyright (c) 2024 Leandro Méndez <leandroa.mendez@gmail.com>
*/
package cmd

import (
	"fmt"

	"github.com/LeanMendez/time-tracker/config"
	"github.com/LeanMendez/time-tracker/utils"
	"github.com/spf13/cobra"
)

var removeCmd = &cobra.Command{
	Use:     "remove [task name or ID]",
	Short:   "remove a task",
	Args:    cobra.ExactArgs(1),
	Aliases: []string{"rm"},
	RunE: func(cmd *cobra.Command, args []string) error {
		taskName := args[0]

		taskManager, err := utils.NewTaskManager(config.ConfigFile)
		if err != nil {
			return err
		}

		tasks, err := taskManager.LoadTasks()
		if err != nil {
			return fmt.Errorf("no tasks found: %w", err)
		}

		_, index, err := taskManager.FindTask(taskName)
		if err != nil {
			return err // Mensaje de error personalizado ya está en FindTask
		}

		tasks = append(tasks[:index], tasks[index+1:]...)

		if err := taskManager.SaveTasks(tasks); err != nil {
			return fmt.Errorf("failed to save tasks: %w", err)
		}

		fmt.Printf("task '%s' removed successfully", taskName)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(removeCmd)
}
