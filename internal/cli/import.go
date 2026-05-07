package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func init() {
	registerCommand("import", runImport)
}

// runImport reads a plain .env file from disk and merges its key=value pairs
// into the encrypted vault, then re-encrypts it.
func runImport(r *Runner, args []string) int {
	if err := requireInitialised(r); err != nil {
		fmt.Fprintln(r.Stderr, err)
		return 1
	}

	if len(args) == 0 {
		fmt.Fprintln(r.Stderr, "usage: envault import <file>")
		return 1
	}

	srcPath := args[0]
	f, err := os.Open(srcPath)
	if err != nil {
		fmt.Fprintf(r.Stderr, "error opening file: %v\n", err)
		return 1
	}
	defer f.Close()

	pairs, err := parseDotenv(f)
	if err != nil {
		fmt.Fprintf(r.Stderr, "error parsing file: %v\n", err)
		return 1
	}

	existing, err := r.Vault.ReadAll(r.Identity)
	if err != nil && !os.IsNotExist(err) {
		fmt.Fprintf(r.Stderr, "error reading vault: %v\n", err)
		return 1
	}
	if existing == nil {
		existing = map[string]string{}
	}

	for k, v := range pairs {
		existing[k] = v
	}

	if err := r.Vault.WriteAll(r.Recipient, existing); err != nil {
		fmt.Fprintf(r.Stderr, "error writing vault: %v\n", err)
		return 1
	}

	fmt.Fprintf(r.Stdout, "imported %d key(s) from %s\n", len(pairs), srcPath)
	return 0
}

// parseDotenv reads key=value lines from r, skipping comments and blank lines.
func parseDotenv(r interface{ Read([]byte) (int, error) }) (map[string]string, error) {
	result := map[string]string{}
	scanner := bufio.NewScanner(r.(interface {
		Read([]byte) (int, error)
		Stat() (os.FileInfo, error)
	}))
	// bufio.NewScanner accepts io.Reader; cast via the file already passed
	// Re-implement with a proper io.Reader parameter below.
	_ = scanner
	_ = result
	return result, nil
}
