package cli

import (
	"fmt"
	"sort"

	"github.com/urfave/cli/v2"
)

func runFmt(c *cli.Context) error {
	r, err := runnerFromContext(c)
	if err != nil {
		return err
	}

	if err := requireInitialised(r); err != nil {
		return err
	}

	if err := requireVault(r); err != nil {
		return err
	}

	raw, err := r.Vault.Read(r.KeyStore)
	if err != nil {
		return fmt.Errorf("reading vault: %w", err)
	}

	keys := make([]string, 0, len(raw))
	for k := range raw {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// Build an ordered map representation to detect if already sorted.
	alreadySorted := isSorted(raw)

	if c.Bool("check") {
		if !alreadySorted {
			fmt.Fprintln(c.App.Writer, "vault is not formatted (run `envault fmt` to fix)")
			return fmt.Errorf("vault is not formatted")
		}
		fmt.Fprintln(c.App.Writer, "vault is already formatted")
		return nil
	}

	if alreadySorted {
		fmt.Fprintln(c.App.Writer, "vault is already formatted")
		return nil
	}

	// Re-write sorted.
	sorted := make(map[string]string, len(raw))
	for k, v := range raw {
		sorted[k] = v
	}

	if err := r.Vault.Write(r.KeyStore, sorted); err != nil {
		return fmt.Errorf("writing vault: %w", err)
	}

	fmt.Fprintf(c.App.Writer, "formatted %d key(s)\n", len(keys))
	return nil
}

// isSorted returns true when the keys in m are already in lexicographic order
// as they would appear after a sort.Strings pass.
func isSorted(m map[string]string) bool {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	for i := 1; i < len(keys); i++ {
		if keys[i] < keys[i-1] {
			return false
		}
	}
	return true
}
