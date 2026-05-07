package cli

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"
)

func init() {
	registerCommand("copy", runCopy)
}

// runCopy copies the value of a secret to the system clipboard.
func runCopy(r *Runner, args []string) int {
	if len(args) < 1 {
		fmt.Fprintln(r.Stderr, "usage: envault copy <KEY>")
		return 1
	}

	if err := requireInitialised(r); err != nil {
		fmt.Fprintln(r.Stderr, err)
		return 1
	}

	key := args[0]

	envMap, err := r.Vault.Read(r.Identity)
	if err != nil {
		fmt.Fprintf(r.Stderr, "error reading vault: %v\n", err)
		return 1
	}

	val, ok := envMap[key]
	if !ok {
		fmt.Fprintf(r.Stderr, "key %q not found\n", key)
		return 1
	}

	if err := copyToClipboard(val); err != nil {
		fmt.Fprintf(r.Stderr, "clipboard error: %v\n", err)
		return 1
	}

	fmt.Fprintf(r.Stdout, "Copied %s to clipboard.\n", key)
	return 0
}

// copyToClipboard writes s to the system clipboard using platform-specific tools.
func copyToClipboard(s string) error {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("pbcopy")
	case "linux":
		cmd = exec.Command("xclip", "-selection", "clipboard")
	case "windows":
		cmd = exec.Command("clip")
	default:
		return fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}
	cmd.Stdin = strings.NewReader(s)
	return cmd.Run()
}
