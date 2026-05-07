package cli

import (
	"fmt"
	"sort"
	"strings"
)

func init() {
	registerCommand("diff", runDiff)
}

// runDiff decrypts the vault and compares its contents against the current
// .env file in the working directory, printing additions, removals and changes.
func runDiff(r *Runner, args []string) int {
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

	vaultEnv, err := r.Vault.Read(identity)
	if err != nil {
		fmt.Fprintln(r.Stderr, "failed to decrypt vault:", err)
		return 1
	}

	dotenvPath := ".env"
	if len(args) > 0 {
		dotenvPath = args[0]
	}

	diskRaw, err := readFileBytes(dotenvPath)
	if err != nil {
		fmt.Fprintln(r.Stderr, "failed to read .env file:", err)
		return 1
	}

	diskEnv, err := parseDotenv(string(diskRaw))
	if err != nil {
		fmt.Fprintln(r.Stderr, "failed to parse .env file:", err)
		return 1
	}

	diffs := computeDiff(vaultEnv, diskEnv)
	if len(diffs) == 0 {
		fmt.Fprintln(r.Stdout, "no differences")
		return 0
	}

	for _, line := range diffs {
		fmt.Fprintln(r.Stdout, line)
	}
	return 0
}

// computeDiff returns a slice of human-readable diff lines comparing vault
// (left) against disk (right).
func computeDiff(vault, disk map[string]string) []string {
	keys := make(map[string]struct{})
	for k := range vault {
		keys[k] = struct{}{}
	}
	for k := range disk {
		keys[k] = struct{}{}
	}

	sorted := make([]string, 0, len(keys))
	for k := range keys {
		sorted = append(sorted, k)
	}
	sort.Strings(sorted)

	var lines []string
	for _, k := range sorted {
		vVal, inVault := vault[k]
		dVal, inDisk := disk[k]
		switch {
		case inVault && !inDisk:
			lines = append(lines, fmt.Sprintf("- %s=%s", k, maskValue(vVal)))
		case !inVault && inDisk:
			lines = append(lines, fmt.Sprintf("+ %s=%s", k, maskValue(dVal)))
		case vVal != dVal:
			lines = append(lines, fmt.Sprintf("~ %s: %s -> %s", k, maskValue(vVal), maskValue(dVal)))
		}
	}
	return lines
}

func maskValue(v string) string {
	if len(v) == 0 {
		return ""
	}
	return strings.Repeat("*", min(len(v), 8))
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
