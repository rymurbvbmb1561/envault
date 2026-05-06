package keystore

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"filippo.io/age"
)

// KeyStore manages the age identity (private key) for a project.
type KeyStore struct {
	path string
}

// New returns a KeyStore whose key file lives at path.
func New(path string) *KeyStore {
	return &KeyStore{path: path}
}

// Exists reports whether the key file is present on disk.
func (k *KeyStore) Exists() bool {
	_, err := os.Stat(k.path)
	return err == nil
}

// Generate creates a new age X25519 identity and writes it to the key file.
// It returns an error if the file already exists.
func (k *KeyStore) Generate() error {
	if k.Exists() {
		return fmt.Errorf("key already exists at %s", k.path)
	}
	identity, err := age.GenerateX25519Identity()
	if err != nil {
		return fmt.Errorf("generating identity: %w", err)
	}
	if err := os.MkdirAll(filepath.Dir(k.path), 0o700); err != nil {
		return fmt.Errorf("creating key directory: %w", err)
	}
	data := identity.String() + "\n"
	if err := os.WriteFile(k.path, []byte(data), 0o600); err != nil {
		return fmt.Errorf("writing key file: %w", err)
	}
	return nil
}

// Load reads the identity from the key file and returns it.
func (k *KeyStore) Load() (*age.X25519Identity, error) {
	data, err := os.ReadFile(k.path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, fmt.Errorf("key file not found at %s; run 'envault init'", k.path)
		}
		return nil, fmt.Errorf("reading key file: %w", err)
	}
	identities, err := age.ParseIdentities(strings.NewReader(string(data)))
	if err != nil {
		return nil, fmt.Errorf("parsing identity: %w", err)
	}
	if len(identities) == 0 {
		return nil, fmt.Errorf("no identities found in key file")
	}
	id, ok := identities[0].(*age.X25519Identity)
	if !ok {
		return nil, fmt.Errorf("unexpected identity type")
	}
	return id, nil
}

// Recipient loads the identity and returns the corresponding public recipient.
func (k *KeyStore) Recipient() (*age.X25519Recipient, error) {
	id, err := k.Load()
	if err != nil {
		return nil, err
	}
	return id.Recipient(), nil
}
