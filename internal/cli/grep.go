package cli

import (
	"flag"
	"fmt"
	"sort"
	"strings"
)

func init() {
	registerCommand("grep", runGrep)
}

// runGrep searches secret values (and optionally keys) for a substring or pattern.
func runGrep(args []string) int {
	fs := flag.NewFlagSet("grep", flag.ContinueOnError)
	showKeys := fs.Bool("keys", false, "also search key names")
	ignoreCase := fs.Bool("i", false, "case-insensitive match")

	if err := fs.Parse(args); err != nil {
		fmt.Println("usage: envault grep [-keys] [-i] <pattern>")
		return 1
	}

	if fs.NArg() < 1 {
		fmt.Println("error: pattern required")
		fmt.Println("usage: envault grep [-keys] [-i] <pattern>")
		return 1
	}

	pattern := fs.Arg(0)

	st, code := requireStatus()
	if code != 0 {
		return code
	}

	secrets, err := st.runner.Vault.Read()
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return 1
	}

	match := func(s string) bool {
		if *ignoreCase {
			return strings.Contains(strings.ToLower(s), strings.ToLower(pattern))
		}
		return strings.Contains(s, pattern)
	}

	keys := make([]string, 0, len(secrets))
	for k := range secrets {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	found := 0
	for _, k := range keys {
		v := secrets[k]
		keyMatches := *showKeys && match(k)
		valMatches := match(v)
		if keyMatches || valMatches {
			fmt.Printf("%s=%s\n", k, v)
			found++
		}
	}

	if found == 0 {
		fmt.Printf("no matches for %q\n", pattern)
	}

	return 0
}
