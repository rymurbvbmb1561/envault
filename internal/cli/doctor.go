package cli

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"
)

func init() {
	registerCommand("doctor", runDoctor)
}

type checkResult struct {
	name   string
	passed bool
	note   string
}

func runDoctor(r *Runner, args []string, stdout io.Writer) int {
	cfgExists := r.Config.Exists()
	keyExists := r.Keystore.Exists()
	vaultExists := r.Vault.Exists()

	checks := []checkResult{
		{
			name:   "config file present",
			passed: cfgExists,
			note:   ifElse(cfgExists, "", "run 'envault init' to initialise"),
		},
		{
			name:   "key file present",
			passed: keyExists,
			note:   ifElse(keyExists, "", "run 'envault init' to generate a key"),
		},
		{
			name:   "vault file present",
			passed: vaultExists,
			note:   ifElse(vaultExists, "", "run 'envault lock' to create a vault"),
		},
		{
			name:   "editor resolvable",
			passed: resolveEditor() != "",
			note:   ifElse(resolveEditor() != "", "", "set $EDITOR or $VISUAL environment variable"),
		},
		{
			name:   "clipboard tool available",
			passed: clipboardAvailable(),
			note:   ifElse(clipboardAvailable(), "", clipboardHint()),
		},
	}

	allPassed := true
	for _, c := range checks {
		icon := "✓"
		if !c.passed {
			icon = "✗"
			allPassed = false
		}
		line := fmt.Sprintf("  %s %s", icon, c.name)
		if c.note != "" {
			line += fmt.Sprintf(" — %s", c.note)
		}
		fmt.Fprintln(stdout, line)
	}

	fmt.Fprintln(stdout)
	if allPassed {
		fmt.Fprintln(stdout, "All checks passed.")
		return 0
	}
	fmt.Fprintln(stdout, "Some checks failed. See notes above.")
	return 1
}

func ifElse(cond bool, a, b string) string {
	if cond {
		return a
	}
	return b
}

func clipboardAvailable() bool {
	switch runtime.GOOS {
	case "darwin":
		_, err := exec.LookPath("pbcopy")
		return err == nil
	case "linux":
		_, errX := exec.LookPath("xclip")
		_, errW := exec.LookPath("xsel")
		_, errWl := exec.LookPath("wl-copy")
		return errX == nil || errW == nil || errWl == nil
	case "windows":
		return true
	default:
		return false
	}
}

func clipboardHint() string {
	switch runtime.GOOS {
	case "linux":
		return "install xclip, xsel, or wl-clipboard"
	default:
		_ = os.Stderr
		return "clipboard not supported on this platform"
	}
}
