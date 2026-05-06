package config

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

const (
	ConfigFileName = ".envault.json"
	DefaultVaultDir = ".envault"
)

// Config holds project-level envault configuration.
type Config struct {
	VaultDir  string `json:"vault_dir"`
	KeyFile   string `json:"key_file"`
	CreatedAt string `json:"created_at"`
}

// DefaultConfig returns a Config populated with default values.
func DefaultConfig() Config {
	return Config{
		VaultDir: DefaultVaultDir,
		KeyFile:  filepath.Join(DefaultVaultDir, "identity.age"),
	}
}

// Write serialises cfg as JSON to path.
func Write(path string, cfg Config) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return err
	}
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o600)
	if err != nil {
		return err
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	return enc.Encode(cfg)
}

// Load reads and deserialises a Config from path.
func Load(path string) (Config, error) {
	f, err := os.Open(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return Config{}, errors.New("config file not found; run 'envault init' first")
		}
		return Config{}, err
	}
	defer f.Close()
	var cfg Config
	if err := json.NewDecoder(f).Decode(&cfg); err != nil {
		return Config{}, err
	}
	return cfg, nil
}

// Exists reports whether the config file at path exists.
func Exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
