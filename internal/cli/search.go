package cli

import (
	"fmt"
	"sort"
	"strings"
)

func init() {
	registerCommand("search", runSearch)
}

// runSearch filters vault keys (and optionally values) by a query string.
func runSearch(r *Runner, args []string) int {
	if len(args) < 1 {
		fmt.Fprintln(r.Stderr, "usage: envault search <QUERY> [--values]")
		return 1
	}

	if err := requireInitialised(r); err != nil {
		fmt.Fprintln(r.Stderr, err)
		return 1
	}

	query := strings.ToLower(args[0])
	searchValues := len(args) > 1 && args[1] == "--values"

	envMap, err := r.Vault.Read(r.Identity)
	if err != nil {
		fmt.Fprintf(r.Stderr, "error reading vault: %v\n", err)
		return 1
	}

	keys := make([]string, 0, len(envMap))
	for k := range envMap {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	matched := 0
	for _, k := range keys {
		v := envMap[k]
		keyMatch := strings.Contains(strings.ToLower(k), query)
		valMatch := searchValues && strings.Contains(strings.ToLower(v), query)
		if keyMatch || valMatch {
			fmt.Fprintf(r.Stdout, "%s=%s\n", k, v)
			matched++
		}
	}

	if matched == 0 {
		fmt.Fprintf(r.Stdout, "no keys matched %q\n", args[0])
	}
	return 0
}
