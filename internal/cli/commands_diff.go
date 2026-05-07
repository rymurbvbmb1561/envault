package cli

// This file exists so that the diff command's init() function in diff.go is
// included in the build and the command is registered with the CLI runner.
// Additional diff-related wiring (e.g. help text) can be added here.

const diffHelp = `Usage: envault diff [<dotenv-file>]

Compare the contents of the encrypted vault against a plain-text .env file.
Defaults to '.env' in the current directory if no path is supplied.

Output symbols:
  +  key present on disk but not in vault
  -  key present in vault but not on disk
  ~  key present in both but with different values (values are masked)
`
