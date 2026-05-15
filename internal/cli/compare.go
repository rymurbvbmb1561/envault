package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/urfave/cli/v2"
)

func init() {}

func runCompare(c *cli.Context) error {
	r, err := newRunner(c)
	if err != nil {
		return err
	}
	if err := requiresInit(r); err != nil {
		return err
	}

	if c.NArg() < 1 {
		return fmt.Errorf("usage: envault compare <snapshot-a> [snapshot-b]")
	}

	showValues := c.Bool("values")
	keysOnly := c.Bool("keys-only")

	snapshotDir := filepath.Join(r.Config.ProjectRoot, ".envault", "snapshots")

	loadSnapshot := func(name string) (map[string]string, error) {
		path := filepath.Join(snapshotDir, name+".age")
		if _, err := os.Stat(path); os.IsNotExist(err) {
			return nil, fmt.Errorf("snapshot %q not found", name)
		}
		return r.Vault.ReadFrom(path, r.Identity)
	}

	aName := c.Args().Get(0)
	aSecrets, err := loadSnapshot(aName)
	if err != nil {
		return err
	}

	var bSecrets map[string]string
	if c.NArg() >= 2 {
		bName := c.Args().Get(1)
		bSecrets, err = loadSnapshot(bName)
		if err != nil {
			return err
		}
	} else {
		if err := requiresVault(r); err != nil {
			return err
		}
		bSecrets, err = r.Vault.Read(r.Identity)
		if err != nil {
			return err
		}
	}

	printComparison(c, aSecrets, bSecrets, showValues, keysOnly)
	return nil
}

func printComparison(c *cli.Context, a, b map[string]string, showValues, keysOnly bool) {
	keys := unionKeys(a, b)
	sort.Strings(keys)

	changes := 0
	for _, k := range keys {
		av, aOk := a[k]
		bv, bOk := b[k]
		switch {
		case !aOk:
			if showValues {
				fmt.Fprintf(c.App.Writer, "+ %s=%s\n", k, bv)
			} else {
				fmt.Fprintf(c.App.Writer, "+ %s\n", k)
			}
			changes++
		case !bOk:
			if showValues {
				fmt.Fprintf(c.App.Writer, "- %s=%s\n", k, av)
			} else {
				fmt.Fprintf(c.App.Writer, "- %s\n", k)
			}
			changes++
		case !keysOnly && av != bv:
			if showValues {
				fmt.Fprintf(c.App.Writer, "~ %s: %s -> %s\n", k, maskValue(av), maskValue(bv))
			} else {
				fmt.Fprintf(c.App.Writer, "~ %s\n", k)
			}
			changes++
		}
	}

	if changes == 0 {
		fmt.Fprintln(c.App.Writer, "no differences found")
	} else {
		fmt.Fprintf(c.App.Writer, "\n%d difference(s) found\n", changes)
	}
}

func unionKeys(a, b map[string]string) []string {
	seen := make(map[string]struct{})
	for k := range a {
		seen[k] = struct{}{}
	}
	for k := range b {
		seen[k] = struct{}{}
	}
	keys := make([]string, 0, len(seen))
	for k := range seen {
		keys = append(keys, k)
	}
	return keys
}
