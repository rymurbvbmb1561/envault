package cli

import (
	"fmt"
	"sort"
	"strings"
)

func init() {
	registerCommand("tag", runTag)
}

// runTag allows adding, removing, and listing tags on secret keys.
// Usage:
//   envault tag add KEY tag1 tag2
//   envault tag remove KEY tag1
//   envault tag list KEY
func runTag(r *Runner, args []string) int {
	if err := requireInit(r); err != nil {
		fmt.Fprintln(r.Stderr, err)
		return 1
	}
	if len(args) < 2 {
		fmt.Fprintln(r.Stderr, "usage: envault tag <add|remove|list> KEY [tags...]")
		return 1
	}

	subcmd := args[0]
	key := args[1]
	rest := args[2:]

	env, err := r.Vault.Read(r.Identity)
	if err != nil {
		fmt.Fprintln(r.Stderr, "error reading vault:", err)
		return 1
	}

	if _, exists := env[key]; !exists && subcmd != "list" {
		fmt.Fprintf(r.Stderr, "key %q not found in vault\n", key)
		return 1
	}

	tagKey := "__tags__" + key

	switch subcmd {
	case "add":
		if len(rest) == 0 {
			fmt.Fprintln(r.Stderr, "provide at least one tag to add")
			return 1
		}
		existing := parseTags(env[tagKey])
		merged := mergeTags(existing, rest)
		env[tagKey] = strings.Join(merged, ",")
		if err := r.Vault.Write(r.Recipient, env); err != nil {
			fmt.Fprintln(r.Stderr, "error writing vault:", err)
			return 1
		}
		fmt.Fprintf(r.Stdout, "tagged %s: %s\n", key, env[tagKey])

	case "remove":
		if len(rest) == 0 {
			fmt.Fprintln(r.Stderr, "provide at least one tag to remove")
			return 1
		}
		existing := parseTags(env[tagKey])
		updated := removeTags(existing, rest)
		if len(updated) == 0 {
			delete(env, tagKey)
		} else {
			env[tagKey] = strings.Join(updated, ",")
		}
		if err := r.Vault.Write(r.Recipient, env); err != nil {
			fmt.Fprintln(r.Stderr, "error writing vault:", err)
			return 1
		}
		fmt.Fprintf(r.Stdout, "updated tags for %s\n", key)

	case "list":
		tags := parseTags(env[tagKey])
		if len(tags) == 0 {
			fmt.Fprintf(r.Stdout, "no tags for %s\n", key)
		} else {
			fmt.Fprintf(r.Stdout, "%s: %s\n", key, strings.Join(tags, ", "))
		}

	default:
		fmt.Fprintf(r.Stderr, "unknown subcommand %q — use add, remove, or list\n", subcmd)
		return 1
	}

	return 0
}

func parseTags(raw string) []string {
	if raw == "" {
		return nil
	}
	parts := strings.Split(raw, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}

func mergeTags(existing, incoming []string) []string {
	seen := make(map[string]struct{}, len(existing))
	for _, t := range existing {
		seen[t] = struct{}{}
	}
	for _, t := range incoming {
		seen[t] = struct{}{}
	}
	out := make([]string, 0, len(seen))
	for t := range seen {
		out = append(out, t)
	}
	sort.Strings(out)
	return out
}

func removeTags(existing, toRemove []string) []string {
	remove := make(map[string]struct{}, len(toRemove))
	for _, t := range toRemove {
		remove[t] = struct{}{}
	}
	out := make([]string, 0, len(existing))
	for _, t := range existing {
		if _, skip := remove[t]; !skip {
			out = append(out, t)
		}
	}
	return out
}
