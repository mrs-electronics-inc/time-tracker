package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
	"time-tracker/config"
)

var editCmd = &cobra.Command{
	Use:     "edit",
	Short:   "Edit the data file directly",
	Long:    `Open the data.json file in your default editor (EDITOR environment variable).`,
	Aliases: []string{"e"},
	RunE: func(cmd *cobra.Command, args []string) error {
		dataFilePath := config.DataFilePath()

		// Get the editor from the EDITOR environment variable, fallback to nano
		editor := os.Getenv("EDITOR")
		if editor == "" {
			editor = "nano"
		}

		// Create and run the editor command
		editorCmd := exec.Command(editor, dataFilePath)
		editorCmd.Stdin = os.Stdin
		editorCmd.Stdout = os.Stdout
		editorCmd.Stderr = os.Stderr

		if err := editorCmd.Run(); err != nil {
			return fmt.Errorf("failed to open editor: %w", err)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(editCmd)
}
