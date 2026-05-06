package cli

import (
	"fmt"
	"os"
	"os/exec"
)

// defaultEditors is the list of editors to try if EDITOR is not set.
var defaultEditors = []string{"nano", "vi", "vim"}

// resolveEditor returns the editor to use, preferring the EDITOR env var.
func resolveEditor() (string, error) {
	if e := os.Getenv("EDITOR"); e != "" {
		return e, nil
	}
	for _, candidate := range defaultEditors {
		if path, err := exec.LookPath(candidate); err == nil {
			return path, nil
		}
	}
	return "", fmt.Errorf("no suitable editor found; set the EDITOR environment variable")
}

// openInEditor opens the given file path in the resolved editor and waits for
// the user to finish editing. It returns any error from launching or waiting.
func openInEditor(filePath string) error {
	editor, err := resolveEditor()
	if err != nil {
		return err
	}
	cmd := exec.Command(editor, filePath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// Edit decrypts the vault to a temporary file, opens it in the user's editor,
// then re-encrypts the (possibly modified) file back into the vault.
func (r *Runner) Edit() error {
	if err := r.requireInitialised(); err != nil {
		return err
	}

	if !r.vault.Exists() {
		return fmt.Errorf("vault does not exist; run 'envault lock' first to create it")
	}

	// Decrypt to a temp file.
	tmp, err := os.CreateTemp("", "envault-edit-*.env")
	if err != nil {
		return fmt.Errorf("creating temp file: %w", err)
	}
	tmpPath := tmp.Name()
	tmp.Close()
	defer os.Remove(tmpPath)

	identity, err := r.keystore.Load()
	if err != nil {
		return fmt.Errorf("loading key: %w", err)
	}

	if err := r.vault.Read(identity, tmpPath); err != nil {
		return fmt.Errorf("decrypting vault: %w", err)
	}

	// Open editor and wait.
	if err := openInEditor(tmpPath); err != nil {
		return fmt.Errorf("editor exited with error: %w", err)
	}

	// Re-encrypt the (possibly modified) temp file.
	recipient, err := r.keystore.Recipient()
	if err != nil {
		return fmt.Errorf("loading recipient: %w", err)
	}

	if err := r.vault.Write(recipient, tmpPath); err != nil {
		return fmt.Errorf("re-encrypting vault: %w", err)
	}

	fmt.Fprintln(r.stdout, "vault updated")
	return nil
}
