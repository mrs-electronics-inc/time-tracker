/*
Copyright (c) 2024 Leandro Méndez <leandroa.mendez@gmail.com>
*/
package config

import (
	"os"
	"path/filepath"
)

type Config struct {
	StoragePath string `JSON:"storagePath"`
}

var ConfigFile = filepath.Join(os.Getenv("HOME"), ".time-tracker-config.json")
