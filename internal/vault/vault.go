package vault

import (
	"fmt"
	"os"
	"path/filepath"
)

const encryptedExt = ".vault"

// Vault manages encrypted .env files within a project directory.
type Vault struct {
	dir string
}

// New creates a Vault rooted at the given directory.
func New(dir string) (*Vault, error) {
	info, err := os.Stat(dir)
	if err != nil {
		return nil, fmt.Errorf("vault: directory %q not accessible: %w", dir, err)
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("vault: %q is not a directory", dir)
	}
	return &Vault{dir: dir}, nil
}

// VaultPath returns the path where an encrypted version of envFile will be stored.
func (v *Vault) VaultPath(envFile string) string {
	base := filepath.Base(envFile)
	return filepath.Join(v.dir, base+encryptedExt)
}

// Exists reports whether the encrypted vault file for envFile already exists.
func (v *Vault) Exists(envFile string) bool {
	_, err := os.Stat(v.VaultPath(envFile))
	return err == nil
}

// Write stores ciphertext into the vault file for envFile, creating or
// overwriting it.
func (v *Vault) Write(envFile string, ciphertext []byte) error {
	dest := v.VaultPath(envFile)
	if err := os.WriteFile(dest, ciphertext, 0600); err != nil {
		return fmt.Errorf("vault: write %q: %w", dest, err)
	}
	return nil
}

// Read returns the raw ciphertext stored in the vault file for envFile.
func (v *Vault) Read(envFile string) ([]byte, error) {
	src := v.VaultPath(envFile)
	data, err := os.ReadFile(src)
	if err != nil {
		return nil, fmt.Errorf("vault: read %q: %w", src, err)
	}
	return data, nil
}

// Delete removes the vault file for envFile.
func (v *Vault) Delete(envFile string) error {
	dest := v.VaultPath(envFile)
	if err := os.Remove(dest); err != nil {
		return fmt.Errorf("vault: delete %q: %w", dest, err)
	}
	return nil
}
