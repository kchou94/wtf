package cfg

import (
	"errors"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
)

const (
	// WtfConfigDirV1 defines the path to the first version of configuration. Do not use this
	WtfConfigDirV1 = "~/.wtf/"

	// WtfConfigDirV2 defines the path to the second version of the configuration. Use this.
	WtfConfigDirV2 = "~/.config/wtf/"
)

// Initialize takes care of settings up the initial state of WTF configuration
// It ensures necessary directories and files exist
func Initialize(hasCustom bool) {
	if !hasCustom {
		migrateOldConfig()
	}
}

// WtfConfigDir returns the absolute path to the configuration directory
func WtfConfigDir() (string, error) {
	configDir := os.Getenv("XDG_CONFIG_HOME")
	if configDir == "" {
		configDir = WtfConfigDirV2
	} else {
		configDir += "/wtf/"
	}
	configDir, err := expandHomeDir(configDir)
	if err != nil {
		return "", err
	}

	return configDir, nil
}

// Expand expands the path to include the home directory if the path is prefixed with `~`. If it isn't prefixed with `~`, the path is returned as-is.
func expandHomeDir(path string) (string, error) {
	if path == "" {
		return path, nil
	}

	if path[0] != '~' {
		return path, nil
	}

	if len(path) > 1 && path[1] != '/' && path[1] != '\\' {
		return "", errors.New("cannot expand user-specific home dir")
	}

	dir, err := home()
	if err != nil {
		return "", err
	}

	return filepath.Join(dir, path[1:]), nil
}

// Dir returns the home directory for the executing user.
// An error is returned if a home directory cannot be detected.
func home() (string, error) {
	currentUser, err := user.Current()
	if err != nil {
		return "", err
	}
	if currentUser.HomeDir == "" {
		return "", errors.New("cannot find user-specific home dir")
	}

	return currentUser.HomeDir, nil
}

// migrateOldConfig copies any existing configuration from the old location
// to the new, XDG-compatible location
func migrateOldConfig() {
	srcDir, _ := expandHomeDir(WtfConfigDirV1)
	destDir, _ := WtfConfigDir()

	// If the old config directory doesn't exist, do not move
	if _, err := os.Stat(srcDir); os.IsNotExist(err) {
		return
	}

	// If the new config directory already exists, do not move
	if _, err := os.Stat(destDir); err == nil {
		return
	}

	// Time to move
	err := Copy(srcDir, destDir)
	if err != nil {
		panic(err)
	}

	// Delete the old directory if the new one exists
	if _, err := os.Stat(destDir); err == nil {
		err := os.RemoveAll(srcDir)
		if err != nil {
			fmt.Println(err)
		}
	}
}
