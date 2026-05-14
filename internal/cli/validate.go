package cli

import (
	"fmt"
	"regexp"
	"strings"
)

var validKeyPattern = regexp.MustCompile(`^[A-Z][A-Z0-9_]*$`)

type ValidationIssue struct {
	Key     string
	Message string
}

func validateSecrets(secrets map[string]string) []ValidationIssue {
	var issues []ValidationIssue

	for key, value := range secrets {
		if !validKeyPattern.MatchString(key) {
			issues = append(issues, ValidationIssue{
				Key:     key,
				Message: "key must be uppercase alphanumeric with underscores and start with a letter",
			})
		}
		if strings.TrimSpace(value) == "" {
			issues = append(issues, ValidationIssue{
				Key:     key,
				Message: "value is empty or whitespace-only",
			})
		}
		if strings.Contains(value, "\n") {
			issues = append(issues, ValidationIssue{
				Key:     key,
				Message: "value contains newline characters",
			})
		}
	}

	return issues
}

func init() {
	registerCommand("validate", "validate secret keys and values against naming rules", runValidate)
}

func runValidate(r *Runner, args []string) int {
	if err := requireInitialised(r); err != nil {
		fmt.Fprintln(r.Stderr, err)
		return 1
	}
	if err := requireVault(r); err != nil {
		fmt.Fprintln(r.Stderr, err)
		return 1
	}

	secrets, err := r.Vault.Read(r.KeyStore)
	if err != nil {
		fmt.Fprintf(r.Stderr, "error reading vault: %v\n", err)
		return 1
	}

	issues := validateSecrets(secrets)
	if len(issues) == 0 {
		fmt.Fprintln(r.Stdout, "✓ all secrets passed validation")
		return 0
	}

	for _, issue := range issues {
		fmt.Fprintf(r.Stdout, "✗ [%s] %s\n", issue.Key, issue.Message)
	}
	fmt.Fprintf(r.Stdout, "\n%d issue(s) found\n", len(issues))
	return 1
}
