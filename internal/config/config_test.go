package config_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/nicholasgasior/envault/internal/config"
)

func tempConfigPath(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	return filepath.Join(dir, config.ConfigFileName)
}

func TestWriteAndLoad(t *testing.T) {
	path := tempConfigPath(t)
	cfg := config.Config{
		VaultDir:  ".envault",
		KeyFile:   ".envault/identity.age",
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
	}
	if err := config.Write(path, cfg); err != nil {
		t.Fatalf("Write: %v", err)
	}
	got, err := config.Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if got.VaultDir != cfg.VaultDir {
		t.Errorf("VaultDir: got %q, want %q", got.VaultDir, cfg.VaultDir)
	}
	if got.KeyFile != cfg.KeyFile {
		t.Errorf("KeyFile: got %q, want %q", got.KeyFile, cfg.KeyFile)
	}
}

func TestLoadMissingFileReturnsError(t *testing.T) {
	_, err := config.Load("/nonexistent/.envault.json")
	if err == nil {
		t.Fatal("expected error for missing config file, got nil")
	}
}

func TestExistsReturnsFalseBeforeWrite(t *testing.T) {
	path := tempConfigPath(t)
	if config.Exists(path) {
		t.Error("Exists should return false before Write")
	}
}

func TestExistsReturnsTrueAfterWrite(t *testing.T) {
	path := tempConfigPath(t)
	if err := config.Write(path, config.DefaultConfig()); err != nil {
		t.Fatalf("Write: %v", err)
	}
	if !config.Exists(path) {
		t.Error("Exists should return true after Write")
	}
}

func TestDefaultConfig(t *testing.T) {
	cfg := config.DefaultConfig()
	if cfg.VaultDir == "" {
		t.Error("DefaultConfig VaultDir should not be empty")
	}
	if cfg.KeyFile == "" {
		t.Error("DefaultConfig KeyFile should not be empty")
	}
}

func TestWriteCreatesParentDirectories(t *testing.T) {
	base := t.TempDir()
	path := filepath.Join(base, "nested", "dir", config.ConfigFileName)
	if err := config.Write(path, config.DefaultConfig()); err != nil {
		t.Fatalf("Write with nested path: %v", err)
	}
	if _, err := os.Stat(path); err != nil {
		t.Errorf("expected config file to exist at %s: %v", path, err)
	}
}
