package cli

import (
	"fmt"
	"io"
	"os"

	"github.com/user/envault/internal/config"
	"github.com/user/envault/internal/keystore"
	"github.com/user/envault/internal/vault"
)

// Runner holds the dependencies needed to execute CLI commands.
type Runner struct {
	keys  *keystore.KeyStore
	vault *vault.Vault
	cfg   *config.Config
	out   io.Writer
}

// New creates a Runner rooted at the given directory.
func New(dir string) (*Runner, error) {
	cfg, err := config.LoadFromCwd()
	if err != nil {
		cfg = config.DefaultConfig(dir)
	}

	ks := keystore.New(cfg.KeyPath)
	v := vault.New(cfg.VaultPath)

	return &Runner{
		keys:  ks,
		vault: v,
		cfg:   cfg,
		out:   os.Stdout,
	}, nil
}

// SetOutput redirects command output (useful in tests).
func (r *Runner) SetOutput(w io.Writer) {
	r.out = w
}

// printf is a small helper that writes formatted output to r.out.
func (r *Runner) printf(format string, args ...any) {
	fmt.Fprintf(r.out, format, args...)
}
