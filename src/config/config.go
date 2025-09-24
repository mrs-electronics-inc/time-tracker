package config

import (
	"os"
	"path/filepath"
)

type Config struct {
	// Config is now empty since we don't need to store anything
}

var ConfigFile = filepath.Join(os.Getenv("HOME"), ".config/time-tracker/config.json")
