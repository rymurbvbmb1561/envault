package cli

import (
	"fmt"
	"sort"

	"github.com/urfave/cli/v2"
)

// vaultPathFromContext returns the resolved vault path for the current project.
func vaultPathFromContext(c *cli.Context) string {
	r := runnerFromContext(c)
	return r.Vault.Path()
}

// loadSecretsFromContext decrypts the vault and returns key-value pairs sorted
// by key for deterministic output.
func loadSecretsFromContext(c *cli.Context) ([][2]string, error) {
	r := runnerFromContext(c)
	data, err := r.Vault.Read(r.Key)
	if err != nil {
		return nil, fmt.Errorf("read vault: %w", err)
	}
	keys := make([]string, 0, len(data))
	for k := range data {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	pairs := make([][2]string, 0, len(keys))
	for _, k := range keys {
		pairs = append(pairs, [2]string{k, data[k]})
	}
	return pairs, nil
}
