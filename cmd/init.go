/*
Copyright (c) 2024 Leandro Méndez <leandroa.mendez@gmail.com>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/LeanMendez/time-tracker/config"
	"github.com/LeanMendez/time-tracker/models"
	"github.com/LeanMendez/time-tracker/utils"
	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init [path]",
	Short: "Initialize time-tracker with a storage path",
	Long:  `Initializes the application by defining a path to the storage file where the data of the tasks created in the application will be saved.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		path := args[0]

		// Create the storage path
		if err := os.MkdirAll(path, 0755); err != nil {
			return fmt.Errorf("failed to create directory: %w", err)
		}

		// Save the storage path to the configuration file
		configLocal := config.Config{
			StoragePath: path,
		}
		configData, err := json.Marshal(configLocal)
		if err != nil {
			return fmt.Errorf("failed to marshal config: %w", err)
		}

		if err := os.WriteFile(config.ConfigFile, configData, 0644); err != nil {
			return fmt.Errorf("failed to write config file: %w", err)
		}

		// This will initialize an empty tasks.json
		taskManager, err := utils.NewTaskManager(config.ConfigFile)
		if err != nil {
			return fmt.Errorf("failed to create tasks.json: %w", err)
		}
		if err := taskManager.SaveTasks([]models.Task{}); err != nil {
			return fmt.Errorf("failed to create tasks.json: %w", err)
		}

		fmt.Printf("Initialized time-tracker with storage path: %s \n", path)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
