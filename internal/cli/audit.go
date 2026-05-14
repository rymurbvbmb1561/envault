package cli

import (
	"fmt"
	"sort"
	"strings"

	"github.com/nicholasgasior/envault/internal/config"
)

func init() {
	c := commandRegistry
	c.Register("audit", runAudit)
}

// runAudit inspects the decrypted vault and reports potential issues such as
// empty values, keys with common placeholder values, or duplicate values.
func runAudit(r *Runner, args []string) int {
	cfg, err := config.LoadFromCwd()
	if err != nil {
		fmt.Fprintln(r.Stderr, errNotInitialised)
		return 1
	}

	if !r.Vault.Exists() {
		fmt.Fprintln(r.Stderr, "vault not found: run 'envault lock' to create one")
		return 1
	}

	identity, err := r.KeyStore.Load()
	if err != nil {
		fmt.Fprintf(r.Stderr, "failed to load key: %v\n", err)
		return 1
	}

	data, err := r.Vault.Read(identity)
	if err != nil {
		fmt.Fprintf(r.Stderr, "failed to read vault: %v\n", err)
		return 1
	}

	_ = cfg

	placeholders := []string{"changeme", "todo", "fixme", "xxx", "placeholder", "secret", "password", "example"}

	type issue struct {
		key    string
		reason string
	}

	var issues []issue
	valueCount := make(map[string][]string)

	for k, v := range data {
		if strings.TrimSpace(v) == "" {
			issues = append(issues, issue{k, "empty value"})
			continue
		}
		for _, p := range placeholders {
			if strings.EqualFold(v, p) {
				issues = append(issues, issue{k, fmt.Sprintf("placeholder value %q", v)})
				break
			}
		}
		valueCount[v] = append(valueCount[v], k)
	}

	for v, keys := range valueCount {
		if len(keys) > 1 {
			sort.Strings(keys)
			issues = append(issues, issue{strings.Join(keys, ", "), fmt.Sprintf("duplicate value %q", maskValue(v))})
		}
	}

	if len(issues) == 0 {
		fmt.Fprintln(r.Stdout, "no issues found")
		return 0
	}

	sort.Slice(issues, func(i, j int) bool { return issues[i].key < issues[j].key })
	fmt.Fprintf(r.Stdout, "found %d issue(s):\n", len(issues))
	for _, iss := range issues {
		fmt.Fprintf(r.Stdout, "  %-30s %s\n", iss.key, iss.reason)
	}
	return 1
}
