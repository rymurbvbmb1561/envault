package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/user/envault/internal/config"
	"github.com/user/envault/internal/vault"
)

func init() {
	registerCommand("snapshot-restore", runSnapshotRestore)
	registerCommand("snapshot-list", runSnapshotList)
}

// runSnapshotList prints all available snapshots for the current project.
func runSnapshotList(r *Runner, args []string) int {
	cfg, err := config.LoadFromCwd()
	if err != nil {
		fmt.Fprintln(r.Stderr, errNotInitialised)
		return 1
	}

	snaps, err := listSnapshots(cfg.VaultPath)
	if err != nil || len(snaps) == 0 {
		fmt.Fprintln(r.Stdout, "no snapshots found")
		return 0
	}

	for _, s := range snaps {
		fmt.Fprintln(r.Stdout, filepath.Base(s))
	}
	return 0
}

// runSnapshotRestore re-encrypts a named snapshot back into the vault.
// Usage: envault snapshot-restore <timestamp>
func runSnapshotRestore(r *Runner, args []string) int {
	if len(args) == 0 {
		fmt.Fprintln(r.Stderr, "usage: envault snapshot-restore <timestamp>")
		return 1
	}

	cfg, err := config.LoadFromCwd()
	if err != nil {
		fmt.Fprintln(r.Stderr, errNotInitialised)
		return 1
	}

	snapshotDir := filepath.Join(filepath.Dir(cfg.VaultPath), "snapshots")
	name := args[0]
	if !strings.HasSuffix(name, ".env") {
		name += ".env"
	}
	snapshotPath := filepath.Join(snapshotDir, name)

	data, err := os.ReadFile(snapshotPath)
	if err != nil {
		fmt.Fprintln(r.Stderr, "error: snapshot not found:", args[0])
		return 1
	}

	ks := keystoreFromConfig(cfg)
	id, err := ks.Load()
	if err != nil {
		fmt.Fprintln(r.Stderr, "error: could not load key:", err)
		return 1
	}

	v := vault.New(cfg.VaultPath)
	if err := v.Write(id, data); err != nil {
		fmt.Fprintln(r.Stderr, "error: could not write vault:", err)
		return 1
	}

	fmt.Fprintln(r.Stdout, "vault restored from snapshot:", args[0])
	return 0
}

func listSnapshots(vaultPath string) ([]string, error) {
	dir := filepath.Join(filepath.Dir(vaultPath), "snapshots")
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	var out []string
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ".env") {
			out = append(out, filepath.Join(dir, e.Name()))
		}
	}
	sort.Strings(out)
	return out, nil
}
