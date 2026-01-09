package cmd

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	"time-tracker/cmd/headless"
	"time-tracker/cmd/tui"
	"time-tracker/config"
	"time-tracker/utils"
)

var rootCmd = &cobra.Command{
	Use:   "time-tracker",
	Short: "A simple time tracker",
	Long: `A simple time tracker TUI application.
See https://github.com/mrs-electronics-inc/time-tracker for more details.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// If no subcommand provided, launch the TUI
		if len(args) == 0 {
			storage, err := utils.NewFileStorage(config.DataFilePath())
			if err != nil {
				return fmt.Errorf("failed to initialize storage: %w", err)
			}
			taskManager := utils.NewTaskManager(storage)

			model := tui.NewModel(storage, taskManager)
			if err := model.LoadEntries(); err != nil {
				return fmt.Errorf("failed to load entries: %w", err)
			}

			p := tea.NewProgram(model, tea.WithAltScreen())
			_, err = p.Run()
			return err
		}
		return cmd.Help()
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(headless.HeadlessCmd)
}
