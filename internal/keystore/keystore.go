package keystore

import (
	"errors"
	"os"
	"path/filepath"

	"filippo.io/age"
)

const (
	defaultKeyDir  = ".envault"
	defaultKeyFile = "identity.age"
)

// KeyStore manages age identity keys on disk.
type KeyStore struct {
	KeyDir string
}

// New returns a KeyStore rooted at the user's home directory.
func New() (*KeyStore, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	return &KeyStore{KeyDir: filepath.Join(home, defaultKeyDir)}, nil
}

// KeyPath returns the full path to the identity file.
func (ks *KeyStore) KeyPath() string {
	return filepath.Join(ks.KeyDir, defaultKeyFile)
}

// Exists reports whether an identity file already exists.
func (ks *KeyStore) Exists() bool {
	_, err := os.Stat(ks.KeyPath())
	return err == nil
}

// Generate creates a new age X25519 identity and writes it to disk.
// Returns an error if a key already exists.
func (ks *KeyStore) Generate() (*age.X25519Identity, error) {
	if ks.Exists() {
		return nil, errors.New("identity already exists; use Load() or delete it first")
	}
	if err := os.MkdirAll(ks.KeyDir, 0700); err != nil {
		return nil, err
	}
	id, err := age.GenerateX25519Identity()
	if err != nil {
		return nil, err
	}
	f, err := os.OpenFile(ks.KeyPath(), os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	_, err = f.WriteString(id.String() + "\n")
	return id, err
}

// Load reads and parses the stored age identity.
func (ks *KeyStore) Load() (*age.X25519Identity, error) {
	data, err := os.ReadFile(ks.KeyPath())
	if err != nil {
		return nil, err
	}
	identities, err := age.ParseIdentities(strings.NewReader(string(data)))
	if err != nil {
		return nil, err
	}
	if len(identities) == 0 {
		return nil, errors.New("no identities found in key file")
	}
	id, ok := identities[0].(*age.X25519Identity)
	if !ok {
		return nil, errors.New("unexpected identity type")
	}
	return id, nil
}
