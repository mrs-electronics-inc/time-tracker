package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"time-tracker/config"
	"time-tracker/models"
	"time-tracker/utils"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize time-tracker",
	Long:  `Initializes the application by creating the configuration and storage files in the user's config directory.`,
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Get the storage path (same directory as config file)
		storagePath := filepath.Dir(config.ConfigFile)

		// Create the storage path
		if err := os.MkdirAll(storagePath, 0755); err != nil {
			return fmt.Errorf("failed to create directory: %w", err)
		}

		// Save the storage path to the configuration file
		configLocal := config.Config{
			StoragePath: storagePath,
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

		fmt.Printf("Initialized time-tracker with storage path: %s \n", storagePath)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
