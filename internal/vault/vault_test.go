package vault_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/envault/internal/vault"
)

func tempVault(t *testing.T) (*vault.Vault, string) {
	t.Helper()
	dir := t.TempDir()
	v, err := vault.New(dir)
	if err != nil {
		t.Fatalf("vault.New: %v", err)
	}
	return v, dir
}

func TestVaultPathHasExtension(t *testing.T) {
	v, _ := tempVault(t)
	got := v.VaultPath(".env")
	if filepath.Ext(got) != ".vault" {
		t.Errorf("expected .vault extension, got %q", got)
	}
}

func TestExistsReturnsFalseInitially(t *testing.T) {
	v, _ := tempVault(t)
	if v.Exists(".env") {
		t.Error("expected Exists to return false before any write")
	}
}

func TestWriteAndRead(t *testing.T) {
	v, _ := tempVault(t)
	data := []byte("secret ciphertext")
	if err := v.Write(".env", data); err != nil {
		t.Fatalf("Write: %v", err)
	}
	got, err := v.Read(".env")
	if err != nil {
		t.Fatalf("Read: %v", err)
	}
	if string(got) != string(data) {
		t.Errorf("round-trip mismatch: got %q, want %q", got, data)
	}
}

func TestExistsAfterWrite(t *testing.T) {
	v, _ := tempVault(t)
	_ = v.Write(".env", []byte("data"))
	if !v.Exists(".env") {
		t.Error("expected Exists to return true after write")
	}
}

func TestDelete(t *testing.T) {
	v, _ := tempVault(t)
	_ = v.Write(".env", []byte("data"))
	if err := v.Delete(".env"); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	if v.Exists(".env") {
		t.Error("expected Exists to return false after delete")
	}
}

func TestNewFailsOnFile(t *testing.T) {
	f, err := os.CreateTemp("", "envault-*")
	if err != nil {
		t.Fatal(err)
	}
	f.Close()
	defer os.Remove(f.Name())

	_, err = vault.New(f.Name())
	if err == nil {
		t.Error("expected error when passing a file path to vault.New")
	}
}
