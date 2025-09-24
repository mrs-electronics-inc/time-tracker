package config

import (
	"os"
	"path/filepath"
)

var ConfigPath = func() string {
	configDir, err := os.UserConfigDir()
	if err != nil {
		// Fallback to home directory if UserConfigDir fails
		homeDir, homeErr := os.UserHomeDir()
		if homeErr != nil {
			// Last resort fallback
			configDir = os.Getenv("HOME")
		} else {
			configDir = homeDir
		}
	}
	return filepath.Join(configDir, "time-tracker")
}()

// DataFilePath returns the path to the data.json file
func DataFilePath() string {
	if path := os.Getenv("TIME_TRACKER_DATA_FILE"); path != "" {
		return path
	}
	return filepath.Join(ConfigPath, "data.json")
}
