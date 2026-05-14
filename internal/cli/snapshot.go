package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/user/envault/internal/config"
	"github.com/user/envault/internal/vault"
)

func init() {
	registerCommand("snapshot", runSnapshot)
}

// runSnapshot saves a timestamped copy of the current decrypted .env
// into a .envault/snapshots/ directory for later restoration.
func runSnapshot(r *Runner, args []string) int {
	cfg, err := config.LoadFromCwd()
	if err != nil {
		fmt.Fprintln(r.Stderr, errNotInitialised)
		return 1
	}

	v := vault.New(cfg.VaultPath)
	if !v.Exists() {
		fmt.Fprintln(r.Stderr, "error: no vault found — run 'envault lock' first")
		return 1
	}

	ks := keystoreFromConfig(cfg)
	id, err := ks.Load()
	if err != nil {
		fmt.Fprintln(r.Stderr, "error: could not load key:", err)
		return 1
	}

	data, err := v.Read(id)
	if err != nil {
		fmt.Fprintln(r.Stderr, "error: could not decrypt vault:", err)
		return 1
	}

	snapshotDir := filepath.Join(filepath.Dir(cfg.VaultPath), "snapshots")
	if err := os.MkdirAll(snapshotDir, 0700); err != nil {
		fmt.Fprintln(r.Stderr, "error: could not create snapshots directory:", err)
		return 1
	}

	ts := time.Now().UTC().Format("20060102T150405Z")
	snapshotPath := filepath.Join(snapshotDir, ts+".env")

	if err := os.WriteFile(snapshotPath, data, 0600); err != nil {
		fmt.Fprintln(r.Stderr, "error: could not write snapshot:", err)
		return 1
	}

	fmt.Fprintln(r.Stdout, "snapshot saved:", snapshotPath)
	return 0
}
