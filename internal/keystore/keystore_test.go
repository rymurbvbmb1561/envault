package keystore

import (
	"os"
	"path/filepath"
	"testing"
)

func tempKeyStore(t *testing.T) *KeyStore {
	t.Helper()
	dir := t.TempDir()
	return &KeyStore{KeyDir: filepath.Join(dir, ".envault")}
}

func TestGenerateCreatesKeyFile(t *testing.T) {
	ks := tempKeyStore(t)
	id, err := ks.Generate()
	if err != nil {
		t.Fatalf("Generate() error: %v", err)
	}
	if id == nil {
		t.Fatal("expected non-nil identity")
	}
	if _, err := os.Stat(ks.KeyPath()); err != nil {
		t.Fatalf("key file not found: %v", err)
	}
}

func TestGenerateFailsIfAlreadyExists(t *testing.T) {
	ks := tempKeyStore(t)
	if _, err := ks.Generate(); err != nil {
		t.Fatalf("first Generate() error: %v", err)
	}
	if _, err := ks.Generate(); err == nil {
		t.Fatal("expected error on second Generate(), got nil")
	}
}

func TestLoadRoundTrip(t *testing.T) {
	ks := tempKeyStore(t)
	original, err := ks.Generate()
	if err != nil {
		t.Fatalf("Generate() error: %v", err)
	}
	loaded, err := ks.Load()
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}
	if original.Recipient().String() != loaded.Recipient().String() {
		t.Errorf("recipient mismatch: got %s, want %s",
			loaded.Recipient().String(), original.Recipient().String())
	}
}

func TestExistsReturnsFalseBeforeGenerate(t *testing.T) {
	ks := tempKeyStore(t)
	if ks.Exists() {
		t.Fatal("expected Exists() to return false before Generate()")
	}
}

func TestKeyFilePermissions(t *testing.T) {
	ks := tempKeyStore(t)
	if _, err := ks.Generate(); err != nil {
		t.Fatalf("Generate() error: %v", err)
	}
	info, err := os.Stat(ks.KeyPath())
	if err != nil {
		t.Fatalf("Stat() error: %v", err)
	}
	if perm := info.Mode().Perm(); perm != 0600 {
		t.Errorf("expected permissions 0600, got %04o", perm)
	}
}
