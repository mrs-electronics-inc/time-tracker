/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/LeanMendez/time-tracker/config"
	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init [path]",
	Short: "Initialize timer-cli with a storage path",
	Long:  `Initializes the application by defining a path to the storage file where the data of the tasks created in the application will be saved.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		path := args[0]

		if err := os.MkdirAll(path, 0755); err != nil {
			return fmt.Errorf("failed to create directory: %w", err)
		}
		
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

		fmt.Printf("Initialized timer-cli with storage path: %s \n", path)
		
		return nil
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
