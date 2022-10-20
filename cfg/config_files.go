package cfg

import (
	"errors"
	"os"
	"os/user"
	"path/filepath"
)

const (
	// WtfConfigDirV2 defines the path to the second version of the configuration. Use this.
	WtfConfigDirV2 = "~/.config/wtf/"
)

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
