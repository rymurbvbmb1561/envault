package config

import (
	"os"
	"path/filepath"
)

// FindProjectRoot walks up from cwd looking for ConfigFileName.
// It returns the directory containing the config file, or an empty string
// if none is found before reaching the filesystem root.
func FindProjectRoot(cwd string) (string, bool) {
	dir := cwd
	for {
		candidate := filepath.Join(dir, ConfigFileName)
		if Exists(candidate) {
			return dir, true
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			// reached filesystem root
			return "", false
		}
		dir = parent
	}
}

// LoadFromCwd locates the config file by walking up from the current working
// directory and loads it. Returns an error if not found or unreadable.
func LoadFromCwd() (Config, string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return Config{}, "", err
	}
	root, ok := FindProjectRoot(cwd)
	if !ok {
		return Config{}, "", &NotFoundError{Cwd: cwd}
	}
	cfgPath := filepath.Join(root, ConfigFileName)
	cfg, err := Load(cfgPath)
	if err != nil {
		return Config{}, "", err
	}
	return cfg, root, nil
}

// NotFoundError is returned when no config file can be located.
type NotFoundError struct {
	Cwd string
}

func (e *NotFoundError) Error() string {
	return "no envault config found in " + e.Cwd + " or any parent directory; run 'envault init'"
}
