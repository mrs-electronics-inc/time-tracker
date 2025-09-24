package config

import (
	"os"
	"path/filepath"
)

type Config struct {
	StoragePath string `JSON:"storagePath"`
}

var ConfigFile = filepath.Join(os.Getenv("HOME"), ".config/time-tracker/config.json")
