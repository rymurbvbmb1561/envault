package cli

import (
	"fmt"
	"strings"
)

func init() {
	registerCommand("rename", runRename)
}

// runRename renames an existing key in the vault to a new name,
// preserving its value. Usage: envault rename OLD_KEY NEW_KEY
func runRename(r *Runner, args []string) int {
	if len(args) < 2 {
		fmt.Fprintln(r.Stderr, "usage: envault rename <old-key> <new-key>")
		return 1
	}

	oldKey := strings.TrimSpace(args[0])
	newKey := strings.TrimSpace(args[1])

	if oldKey == "" || newKey == "" {
		fmt.Fprintln(r.Stderr, "error: key names must not be empty")
		return 1
	}

	if oldKey == newKey {
		fmt.Fprintln(r.Stderr, "error: old and new key names are the same")
		return 1
	}

	if !isValidKey(newKey) {
		fmt.Fprintf(r.Stderr, "error: invalid key name %q\n", newKey)
		return 1
	}

	ok, err := checkInitialised(r)
	if !ok {
		fmt.Fprintln(r.Stderr, err)
		return 1
	}

	env, loadErr := r.Vault.Read(r.Identity)
	if loadErr != nil {
		fmt.Fprintf(r.Stderr, "error: could not read vault: %v\n", loadErr)
		return 1
	}

	value, exists := env[oldKey]
	if !exists {
		fmt.Fprintf(r.Stderr, "error: key %q not found\n", oldKey)
		return 1
	}

	if _, taken := env[newKey]; taken {
		fmt.Fprintf(r.Stderr, "error: key %q already exists; remove it first\n", newKey)
		return 1
	}

	env[newKey] = value
	delete(env, oldKey)

	if writeErr := r.Vault.Write(r.Recipient, env); writeErr != nil {
		fmt.Fprintf(r.Stderr, "error: could not write vault: %v\n", writeErr)
		return 1
	}

	fmt.Fprintf(r.Stdout, "renamed %q → %q\n", oldKey, newKey)
	return 0
}
