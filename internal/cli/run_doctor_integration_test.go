package cli

import (
	"bytes"
	"strings"
	"testing"
)

// TestRunDoctorCommandDispatch verifies that "doctor" is wired into the
// top-level Run dispatcher and returns a sensible exit code.
func TestRunDoctorCommandDispatch(t *testing.T) {
	r := tempRunner(t)

	var out bytes.Buffer
	code := r.Run([]string{"doctor"}, &out)
	output := out.String()

	// Without init the command must exit 1 and mention at least one failed check.
	if code != 1 {
		t.Errorf("expected exit code 1, got %d", code)
	}
	if !strings.Contains(output, "✗") {
		t.Errorf("expected at least one failing check marker, got:\n%s", output)
	}
}

// TestRunDoctorAfterFullSetup verifies that after init + lock the three core
// checks (config, key, vault) all pass.
func TestRunDoctorAfterFullSetup(t *testing.T) {
	r := tempRunner(t)

	var out bytes.Buffer
	if code := r.Run([]string{"init"}, &out); code != 0 {
		t.Fatalf("init failed: %s", out.String())
	}
	out.Reset()
	if code := r.Run([]string{"lock"}, &out); code != 0 {
		t.Fatalf("lock failed: %s", out.String())
	}
	out.Reset()

	r.Run([]string{"doctor"}, &out)
	output := out.String()

	for _, check := range []string{
		"✓ config file present",
		"✓ key file present",
		"✓ vault file present",
	} {
		if !strings.Contains(output, check) {
			t.Errorf("expected %q in output, got:\n%s", check, output)
		}
	}
}
