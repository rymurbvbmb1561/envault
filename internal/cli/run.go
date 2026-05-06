package cli

import (
	"flag"
	"fmt"
	"os"
)

// Run is the single entry-point for the envault binary. It parses top-level
// flags and dispatches to the appropriate Runner method.
func Run(version string) int {
	if len(os.Args) < 2 {
		printUsage()
		return 1
	}

	r, err := New(version)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}

	cmd := os.Args[1]
	args := os.Args[2:]

	switch cmd {
	case "init":
		err = r.Init()

	case "lock":
		err = r.Lock()

	case "unlock":
		err = r.Unlock()

	case "edit":
		err = r.Edit()

	case "export":
		fs := flag.NewFlagSet("export", flag.ContinueOnError)
		format := fs.String("format", "export", "output format: export|dotenv")
		if err2 := fs.Parse(args); err2 != nil {
			return 1
		}
		err = r.Export(*format)

	case "status":
		err = r.Status()

	case "version":
		err = r.Version()

	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n", cmd)
		printUsage()
		return 1
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}
	return 0
}

func printUsage() {
	fmt.Fprintln(os.Stderr, `envault — local secrets manager

Usage:
  envault <command> [flags]

Commands:
  init      Initialise envault in the current project
  lock      Encrypt .env into the vault
  unlock    Decrypt the vault back to .env
  edit      Open the decrypted .env in your editor
  export    Print secrets as shell exports  (--format export|dotenv)
  status    Show vault status
  version   Print version information`)
}
