package cli

import (
	"errors"
	"fmt"
	"strings"
)

// addCmd handles adding or updating a single key=value pair in the vault.
func (r *Runner) addCmd(args []string) error {
	if err := r.status.RequireInitialised(); err != nil {
		return err
	}

	if len(args) < 1 {
		return errors.New("usage: envault add KEY=VALUE")
	}

	pair := args[0]
	idx := strings.IndexByte(pair, '=')
	if idx <= 0 {
		return fmt.Errorf("invalid argument %q: expected KEY=VALUE", pair)
	}

	key := pair[:idx]
	value := pair[idx+1:]

	if strings.TrimSpace(key) == "" {
		return errors.New("key must not be empty")
	}

	recipient, err := r.keystore.Load()
	if err != nil {
		return fmt.Errorf("load key: %w", err)
	}

	env, err := r.vault.ReadEnv(recipient)
	if err != nil && !errors.Is(err, errVaultNotFound) {
		return fmt.Errorf("read vault: %w", err)
	}
	if env == nil {
		env = map[string]string{}
	}

	env[key] = value

	if err := r.vault.WriteEnv(recipient, env); err != nil {
		return fmt.Errorf("write vault: %w", err)
	}

	fmt.Fprintf(r.stdout, "added %s\n", key)
	return nil
}
