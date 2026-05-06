package vault

import (
	"bytes"
	"fmt"
	"io"

	"filippo.io/age"
)

// Encrypt encrypts plaintext using the provided age recipient and returns the
// armored ciphertext.
func Encrypt(plaintext []byte, recipient age.Recipient) ([]byte, error) {
	var buf bytes.Buffer
	w, err := age.Encrypt(&buf, recipient)
	if err != nil {
		return nil, fmt.Errorf("encrypt: initialise writer: %w", err)
	}
	if _, err := w.Write(plaintext); err != nil {
		return nil, fmt.Errorf("encrypt: write plaintext: %w", err)
	}
	if err := w.Close(); err != nil {
		return nil, fmt.Errorf("encrypt: finalise: %w", err)
	}
	return buf.Bytes(), nil
}

// Decrypt decrypts ciphertext using the provided age identity and returns the
// original plaintext.
func Decrypt(ciphertext []byte, identity age.Identity) ([]byte, error) {
	r, err := age.Decrypt(bytes.NewReader(ciphertext), identity)
	if err != nil {
		return nil, fmt.Errorf("decrypt: open: %w", err)
	}
	plaintext, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("decrypt: read: %w", err)
	}
	return plaintext, nil
}
