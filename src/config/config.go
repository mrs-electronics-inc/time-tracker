package config

import (
	"os"
	"path/filepath"
)

type Config struct {
	// Config is now empty since we don't need to store anything
}

var ConfigFile = func() string {
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
	return filepath.Join(configDir, "time-tracker", "config.json")
}()
