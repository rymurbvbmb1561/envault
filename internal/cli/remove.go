package cli

import (
	"errors"
	"fmt"
)

// removeCmd deletes a key from the vault.
func (r *Runner) removeCmd(args []string) error {
	if err := r.status.RequireInitialised(); err != nil {
		return err
	}

	if len(args) < 1 {
		return errors.New("usage: envault remove KEY")
	}

	key := args[0]
	if key == "" {
		return errors.New("key must not be empty")
	}

	recipient, err := r.keystore.Load()
	if err != nil {
		return fmt.Errorf("load key: %w", err)
	}

	env, err := r.vault.ReadEnv(recipient)
	if err != nil {
		return fmt.Errorf("read vault: %w", err)
	}

	if _, ok := env[key]; !ok {
		return fmt.Errorf("key %q not found in vault", key)
	}

	delete(env, key)

	if err := r.vault.WriteEnv(recipient, env); err != nil {
		return fmt.Errorf("write vault: %w", err)
	}

	fmt.Fprintf(r.stdout, "removed %s\n", key)
	return nil
}
