package cli_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/envault/internal/cli"
)

func tempRunner(t *testing.T) (*cli.Runner, string) {
	t.Helper()
	dir := t.TempDir()
	r, err := cli.New(dir)
	if err != nil {
		t.Fatalf("cli.New: %v", err)
	}
	return r, dir
}

func TestInitCreatesKey(t *testing.T) {
	r, _ := tempRunner(t)
	if err := r.Init(); err != nil {
		t.Fatalf("Init: %v", err)
	}
}

func TestInitFailsIfAlreadyInitialised(t *testing.T) {
	r, _ := tempRunner(t)
	if err := r.Init(); err != nil {
		t.Fatalf("first Init: %v", err)
	}
	if err := r.Init(); err == nil {
		t.Fatal("expected error on second Init, got nil")
	}
}

func TestLockAndUnlockRoundTrip(t *testing.T) {
	r, dir := tempRunner(t)
	if err := r.Init(); err != nil {
		t.Fatalf("Init: %v", err)
	}

	envPath := filepath.Join(dir, ".env")
	original := []byte("DB_HOST=localhost\nDB_PORT=5432\n")
	if err := os.WriteFile(envPath, original, 0600); err != nil {
		t.Fatalf("write .env: %v", err)
	}

	if err := r.Lock(envPath); err != nil {
		t.Fatalf("Lock: %v", err)
	}

	// plaintext file should be gone after locking
	if _, err := os.Stat(envPath); !os.IsNotExist(err) {
		t.Fatal("expected .env to be removed after Lock")
	}

	if err := r.Unlock(envPath); err != nil {
		t.Fatalf("Unlock: %v", err)
	}

	got, err := os.ReadFile(envPath)
	if err != nil {
		t.Fatalf("read recovered .env: %v", err)
	}
	if string(got) != string(original) {
		t.Fatalf("content mismatch: got %q, want %q", got, original)
	}
}

func TestUnlockFailsWithoutVault(t *testing.T) {
	r, dir := tempRunner(t)
	if err := r.Init(); err != nil {
		t.Fatalf("Init: %v", err)
	}
	envPath := filepath.Join(dir, ".env")
	if err := r.Unlock(envPath); err == nil {
		t.Fatal("expected error unlocking non-existent vault")
	}
}
