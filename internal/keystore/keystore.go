package keystore

import (
	"errors"
	"fmt"
	"os"

	"filippo.io/age"
)

// KeyStore manages a single age X25519 identity stored in a file.
type KeyStore struct {
	path string
}

// New returns a KeyStore rooted at path.
func New(path string) *KeyStore {
	return &KeyStore{path: path}
}

// Exists reports whether the key file is present on disk.
func (k *KeyStore) Exists() bool {
	_, err := os.Stat(k.path)
	return err == nil
}

// Generate creates a new age identity and writes it to the key file.
// It returns an error if the file already exists.
func (k *KeyStore) Generate() error {
	if k.Exists() {
		return fmt.Errorf("key file already exists at %s", k.path)
	}
	return k.generate()
}

// Rotate generates a brand-new identity and overwrites the existing key file.
func (k *KeyStore) Rotate() error {
	return k.generate()
}

// generate is the shared implementation for Generate and Rotate.
func (k *KeyStore) generate() error {
	identity, err := age.GenerateX25519Identity()
	if err != nil {
		return fmt.Errorf("generating age identity: %w", err)
	}

	f, err := os.OpenFile(k.path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return fmt.Errorf("creating key file: %w", err)
	}
	defer f.Close()

	_, err = fmt.Fprintln(f, identity.String())
	return err
}

// Load reads and parses the age identity from the key file.
func (k *KeyStore) Load() (*age.X25519Identity, error) {
	data, err := os.ReadFile(k.path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, fmt.Errorf("key file not found – run 'envault init' first")
		}
		return nil, fmt.Errorf("reading key file: %w", err)
	}

	identities, err := age.ParseIdentities(strings.NewReader(strings.TrimSpace(string(data))))
	if err != nil {
		return nil, fmt.Errorf("parsing age identity: %w", err)
	}
	if len(identities) == 0 {
		return nil, fmt.Errorf("no identity found in key file")
	}

	id, ok := identities[0].(*age.X25519Identity)
	if !ok {
		return nil, fmt.Errorf("unexpected identity type")
	}
	return id, nil
}
