package cli

import (
	"fmt"
	"os"
	"sort"
	"strings"
)

func init() {
	registerCommand("env", runEnv)
}

// runEnv prints the decrypted secrets as KEY=VALUE pairs suitable for
// evaluation in the current shell (e.g. eval $(envault env)).
// With --export it prefixes each line with "export ".
func runEnv(r *Runner, args []string) int {
	exportPrefix := false
	filtered := args[:0]
	for _, a := range args {
		if a == "--export" || a == "-e" {
			exportPrefix = true
		} else {
			filtered = append(filtered, a)
		}
	}
	args = filtered

	if err := requireInitialised(r); err != nil {
		fmt.Fprintln(r.Stderr, err)
		return 1
	}

	if !r.Vault.Exists() {
		fmt.Fprintln(r.Stderr, "no vault found – run 'envault lock' first")
		return 1
	}

	identity, err := r.KeyStore.Load()
	if err != nil {
		fmt.Fprintln(r.Stderr, "failed to load key:", err)
		return 1
	}

	data, err := r.Vault.Read(identity)
	if err != nil {
		fmt.Fprintln(r.Stderr, "failed to decrypt vault:", err)
		return 1
	}

	// Optional prefix filter: envault env AWS
	prefix := ""
	if len(args) > 0 {
		prefix = strings.ToUpper(args[0])
	}

	keys := make([]string, 0, len(data))
	for k := range data {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		if prefix != "" && !strings.HasPrefix(k, prefix) {
			continue
		}
		v := data[k]
		// Shell-quote the value by wrapping in single quotes and escaping
		// any single quotes inside the value.
		quoted := "'" + strings.ReplaceAll(v, "'", "'\\''")	+ "'"
		if exportPrefix {
			fmt.Fprintf(os.Stdout, "export %s=%s\n", k, quoted)
		} else {
			fmt.Fprintf(os.Stdout, "%s=%s\n", k, quoted)
		}
	}
	return 0
}
