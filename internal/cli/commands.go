package cli

import (
	"fmt"
	"os"

	"github.com/user/envault/internal/keystore"
	"github.com/user/envault/internal/vault"
)

// Runner holds dependencies for CLI command execution.
type Runner struct {
	ks *keystore.KeyStore
	v  *vault.Vault
}

// New creates a Runner rooted at the given directory.
func New(dir string) (*Runner, error) {
	ks, err := keystore.New(dir)
	if err != nil {
		return nil, fmt.Errorf("keystore: %w", err)
	}
	v, err := vault.New(dir)
	if err != nil {
		return nil, fmt.Errorf("vault: %w", err)
	}
	return &Runner{ks: ks, v: v}, nil
}

// Init generates a new age identity for the project.
// Returns an error if a key already exists.
func (r *Runner) Init() error {
	if err := r.ks.Generate(); err != nil {
		return fmt.Errorf("init: %w", err)
	}
	fmt.Println("envault: key generated successfully")
	return nil
}

// Lock encrypts the plaintext .env file into the vault.
// The plaintext file is removed after successful encryption.
func (r *Runner) Lock(envPath string) error {
	identity, err := r.ks.Load()
	if err != nil {
		return fmt.Errorf("lock: %w", err)
	}
	plaintext, err := os.ReadFile(envPath)
	if err != nil {
		return fmt.Errorf("lock: read .env: %w", err)
	}
	if err := r.v.Write(plaintext, identity); err != nil {
		return fmt.Errorf("lock: encrypt: %w", err)
	}
	if err := os.Remove(envPath); err != nil {
		return fmt.Errorf("lock: remove plaintext: %w", err)
	}
	fmt.Printf("envault: %s locked\n", envPath)
	return nil
}

// Unlock decrypts the vault back into a plaintext .env file.
func (r *Runner) Unlock(envPath string) error {
	identity, err := r.ks.Load()
	if err != nil {
		return fmt.Errorf("unlock: %w", err)
	}
	plaintext, err := r.v.Read(identity)
	if err != nil {
		return fmt.Errorf("unlock: decrypt: %w", err)
	}
	if err := os.WriteFile(envPath, plaintext, 0600); err != nil {
		return fmt.Errorf("unlock: write .env: %w", err)
	}
	fmt.Printf("envault: %s unlocked\n", envPath)
	return nil
}
