package cli

import (
	"fmt"
	"sort"
)

// listCmd prints all keys stored in the vault, one per line.
func (r *Runner) listCmd(_ []string) error {
	if err := r.status.RequireInitialised(); err != nil {
		return err
	}

	recipient, err := r.keystore.Load()
	if err != nil {
		return fmt.Errorf("load key: %w", err)
	}

	env, err := r.vault.ReadEnv(recipient)
	if err != nil {
		return fmt.Errorf("read vault: %w", err)
	}

	if len(env) == 0 {
		fmt.Fprintln(r.stdout, "(no secrets stored)")
		return nil
	}

	keys := make([]string, 0, len(env))
	for k := range env {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		fmt.Fprintln(r.stdout, k)
	}
	return nil
}
