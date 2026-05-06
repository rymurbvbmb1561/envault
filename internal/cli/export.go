package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/user/envault/internal/config"
)

// exportFormat controls how secrets are printed to stdout.
type exportFormat int

const (
	formatExport exportFormat = iota // export KEY=VALUE
	formatDotenv                     // KEY=VALUE
)

// Export decrypts the vault and prints its contents to stdout in the
// requested format so the caller can source or eval the output.
func (r *Runner) Export(format string) error {
	if err := r.requireInitialised(); err != nil {
		return err
	}

	cfg, err := config.LoadFromCwd()
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	identity, err := r.ks.Load()
	if err != nil {
		return fmt.Errorf("load key: %w", err)
	}

	if !r.vault.Exists() {
		return fmt.Errorf("no vault found — run 'envault lock' first")
	}

	ciphertext, err := r.vault.Read()
	if err != nil {
		return fmt.Errorf("read vault: %w", err)
	}

	plaintext, err := r.enc.Decrypt(ciphertext, identity)
	if err != nil {
		return fmt.Errorf("decrypt vault: %w", err)
	}

	_ = cfg // reserved for future per-project options

	lines := strings.Split(strings.TrimRight(string(plaintext), "\n"), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		switch strings.ToLower(format) {
		case "dotenv", "dot":
			fmt.Fprintln(os.Stdout, line)
		default:
			fmt.Fprintf(os.Stdout, "export %s\n", line)
		}
	}
	return nil
}
