package cmd

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	"time-tracker/cmd/tui"
	"time-tracker/config"
	"time-tracker/utils"
)

var tuiCmd = &cobra.Command{
	Use:   "tui",
	Short: "Launch the interactive TUI interface",
	Long:  `Launch the interactive Text User Interface for time tracking.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		storage, err := utils.NewFileStorage(config.DataFilePath())
		if err != nil {
			return fmt.Errorf("failed to initialize storage: %w", err)
		}
		taskManager := utils.NewTaskManager(storage)

		model := tui.NewModel(storage, taskManager)
		if err := model.LoadEntries(); err != nil {
			return fmt.Errorf("failed to load entries: %w", err)
		}

		p := tea.NewProgram(model)
		_, err = p.Run()
		return err
	},
}

func init() {
	rootCmd.AddCommand(tuiCmd)
}
