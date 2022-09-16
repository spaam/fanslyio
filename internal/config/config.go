package config

import (
	"os"
	"path/filepath"
	"runtime"
)

func Configfile() string {
	userHome, _ := os.UserHomeDir()
	switch runtime.GOOS {
	case "windows":
		return filepath.Join(os.Getenv("AppData"), "fanslyio/config.json")
	case "darwin":
		return filepath.Join(userHome, "Library/Application Support/fanslyio/config.json")
	default:
		return filepath.Join(userHome, ".config/fanslyio/config.json")
	}
}
