package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/nicholasgasior/envault/internal/config"
)

func TestFindProjectRootFindsConfig(t *testing.T) {
	root := t.TempDir()
	cfgPath := filepath.Join(root, config.ConfigFileName)
	if err := config.Write(cfgPath, config.DefaultConfig()); err != nil {
		t.Fatalf("Write: %v", err)
	}

	// Search from a nested subdirectory.
	sub := filepath.Join(root, "a", "b", "c")
	if err := os.MkdirAll(sub, 0o700); err != nil {
		t.Fatalf("MkdirAll: %v", err)
	}

	got, ok := config.FindProjectRoot(sub)
	if !ok {
		t.Fatal("FindProjectRoot should have found the config")
	}
	if got != root {
		t.Errorf("got root %q, want %q", got, root)
	}
}

func TestFindProjectRootReturnsFalseWhenMissing(t *testing.T) {
	dir := t.TempDir()
	_, ok := config.FindProjectRoot(dir)
	if ok {
		t.Error("FindProjectRoot should return false when no config exists")
	}
}

func TestNotFoundErrorMessage(t *testing.T) {
	err := &config.NotFoundError{Cwd: "/some/path"}
	if err.Error() == "" {
		t.Error("NotFoundError.Error() should not be empty")
	}
}

func TestFindProjectRootUsesNearestAncestor(t *testing.T) {
	// Create two nested configs; the nearest ancestor should win.
	outer := t.TempDir()
	inner := filepath.Join(outer, "inner")
	if err := os.MkdirAll(inner, 0o700); err != nil {
		t.Fatalf("MkdirAll: %v", err)
	}

	for _, dir := range []string{outer, inner} {
		if err := config.Write(filepath.Join(dir, config.ConfigFileName), config.DefaultConfig()); err != nil {
			t.Fatalf("Write: %v", err)
		}
	}

	sub := filepath.Join(inner, "sub")
	if err := os.MkdirAll(sub, 0o700); err != nil {
		t.Fatalf("MkdirAll: %v", err)
	}

	got, ok := config.FindProjectRoot(sub)
	if !ok {
		t.Fatal("expected config to be found")
	}
	if got != inner {
		t.Errorf("got %q, want %q (nearest ancestor)", got, inner)
	}
}
