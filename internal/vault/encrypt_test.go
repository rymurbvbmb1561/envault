package vault_test

import (
	"testing"

	"filippo.io/age"
	"github.com/user/envault/internal/vault"
)

func generateIdentity(t *testing.T) *age.X25519Identity {
	t.Helper()
	id, err := age.GenerateX25519Identity()
	if err != nil {
		t.Fatalf("GenerateX25519Identity: %v", err)
	}
	return id
}

func TestEncryptDecryptRoundTrip(t *testing.T) {
	id := generateIdentity(t)
	plaintext := []byte("DB_PASSWORD=supersecret\nAPI_KEY=abc123")

	ciphertext, err := vault.Encrypt(plaintext, id.Recipient())
	if err != nil {
		t.Fatalf("Encrypt: %v", err)
	}
	if len(ciphertext) == 0 {
		t.Fatal("expected non-empty ciphertext")
	}

	got, err := vault.Decrypt(ciphertext, id)
	if err != nil {
		t.Fatalf("Decrypt: %v", err)
	}
	if string(got) != string(plaintext) {
		t.Errorf("round-trip mismatch: got %q, want %q", got, plaintext)
	}
}

func TestDecryptWithWrongKeyFails(t *testing.T) {
	id := generateIdentity(t)
	wrongID := generateIdentity(t)

	ciphertext, err := vault.Encrypt([]byte("secret"), id.Recipient())
	if err != nil {
		t.Fatalf("Encrypt: %v", err)
	}

	_, err = vault.Decrypt(ciphertext, wrongID)
	if err == nil {
		t.Error("expected decryption with wrong key to fail")
	}
}

func TestEncryptProducesDifferentCiphertexts(t *testing.T) {
	id := generateIdentity(t)
	plaintext := []byte("same plaintext")

	c1, _ := vault.Encrypt(plaintext, id.Recipient())
	c2, _ := vault.Encrypt(plaintext, id.Recipient())

	if string(c1) == string(c2) {
		t.Error("expected different ciphertexts for same plaintext (age uses ephemeral keys)")
	}
}
