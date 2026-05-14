package cli

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/nicholasgasior/envault/internal/cli/runner"
)

var (
	keyPattern     = regexp.MustCompile(`^[A-Z][A-Z0-9_]*$`)
	lowerKeyPattern = regexp.MustCompile(`[a-z]`)
)

type lintIssue struct {
	Key     string
	Message string
}

func init() {
	registerCommand("lint", runLint)
}

func runLint(r *runner.Runner, args []string) int {
	if err := requireInit(r); err != nil {
		fmt.Fprintln(r.Stderr, err)
		return 1
	}
	if err := requireVault(r); err != nil {
		fmt.Fprintln(r.Stderr, err)
		return 1
	}

	data, err := r.Vault.Read(r.Keystore)
	if err != nil {
		fmt.Fprintf(r.Stderr, "error reading vault: %v\n", err)
		return 1
	}

	var issues []lintIssue
	for k := range data {
		if lowerKeyPattern.MatchString(k) {
			issues = append(issues, lintIssue{Key: k, Message: "key contains lowercase letters (convention: UPPER_SNAKE_CASE)"})
		} else if !keyPattern.MatchString(k) {
			issues = append(issues, lintIssue{Key: k, Message: "key contains invalid characters (only A-Z, 0-9, _ allowed, must start with a letter)"})
		}
		if strings.HasPrefix(k, "_") {
			issues = append(issues, lintIssue{Key: k, Message: "key starts with underscore"})
		}
		if strings.HasSuffix(k, "_") {
			issues = append(issues, lintIssue{Key: k, Message: "key ends with underscore"})
		}
	}

	if len(issues) == 0 {
		fmt.Fprintln(r.Stdout, "no lint issues found")
		return 0
	}

	fmt.Fprintf(r.Stdout, "%d issue(s) found:\n", len(issues))
	for _, issue := range issues {
		fmt.Fprintf(r.Stdout, "  [%s] %s\n", issue.Key, issue.Message)
	}
	return 1
}
