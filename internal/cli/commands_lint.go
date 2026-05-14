package cli

// lint command is registered via init() in lint.go.
// This file exists to document the lint subcommand and its usage.
//
// Usage:
//   envault lint
//
// Checks all keys stored in the vault for naming convention issues:
//   - Keys should be UPPER_SNAKE_CASE (A-Z, 0-9, _)
//   - Keys must start with a letter
//   - Keys must not start or end with an underscore
//
// Exits with status 0 if no issues are found, 1 otherwise.
