/*
Copyright (c) 2024 Leandro Méndez <leandroa.mendez@gmail.com>
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "time-tracker",
	Short: "A simple time tracker cli",
	Long: `A simple time tracker cli application. 

Time-tracker is a CLI library that allows you to add tasks to a list and track the time until complete them.
This application generate a JSON file where all the data is stored.
Also you can check the information of your tasks listing them.`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
