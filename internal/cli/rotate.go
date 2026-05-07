package cli

import (
	"fmt"
	"io"
)

// rotate re-encrypts the vault with a freshly generated key, replacing the
// existing key-pair stored in the keystore. The plaintext .env is never
// written to disk during the operation.
func (r *Runner) rotate(out io.Writer) error {
	if err := r.checkInitialised(); err != nil {
		return err
	}

	if !r.vault.Exists() {
		return fmt.Errorf("no vault found – lock your secrets first with 'envault lock'")
	}

	// 1. Decrypt with the current key.
	old, err := r.keystore.Load()
	if err != nil {
		return fmt.Errorf("loading current key: %w", err)
	}

	plaintext, err := r.vault.Read(old)
	if err != nil {
		return fmt.Errorf("decrypting vault: %w", err)
	}

	// 2. Generate a replacement key (overwrites the keystore file).
	if err := r.keystore.Rotate(); err != nil {
		return fmt.Errorf("rotating key: %w", err)
	}

	newIdentity, err := r.keystore.Load()
	if err != nil {
		return fmt.Errorf("loading new key: %w", err)
	}

	// 3. Re-encrypt with the new key.
	if err := r.vault.Write(newIdentity, plaintext); err != nil {
		return fmt.Errorf("re-encrypting vault: %w", err)
	}

	fmt.Fprintln(out, "Key rotated and vault re-encrypted successfully.")
	return nil
}
