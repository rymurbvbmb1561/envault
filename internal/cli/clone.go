package cli

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

func runClone(c *cli.Context) error {
	r, err := runnerFromContext(c)
	if err != nil {
		return err
	}

	if err := requireInitialised(r); err != nil {
		return err
	}

	if c.NArg() < 2 {
		return fmt.Errorf("usage: envault clone <source> <destination>")
	}

	srcKey := c.Args().Get(0)
	dstKey := c.Args().Get(1)

	if srcKey == dstKey {
		return fmt.Errorf("source and destination keys must differ")
	}

	v := r.Vault
	if !v.Exists() {
		return fmt.Errorf("no vault found — run 'envault lock' first")
	}

	ks := r.KeyStore
	identity, err := ks.Load()
	if err != nil {
		return fmt.Errorf("loading key: %w", err)
	}

	secrets, err := v.Read(identity)
	if err != nil {
		return fmt.Errorf("reading vault: %w", err)
	}

	value, ok := secrets[srcKey]
	if !ok {
		return fmt.Errorf("key %q not found", srcKey)
	}

	if _, exists := secrets[dstKey]; exists && !c.Bool("force") {
		return fmt.Errorf("key %q already exists; use --force to overwrite", dstKey)
	}

	secrets[dstKey] = value

	recipient, err := identity.Recipient()
	if err != nil {
		return fmt.Errorf("deriving recipient: %w", err)
	}

	if err := v.Write(secrets, recipient); err != nil {
		return fmt.Errorf("writing vault: %w", err)
	}

	fmt.Fprintf(c.App.Writer, "cloned %s → %s\n", srcKey, dstKey)
	return nil
}
