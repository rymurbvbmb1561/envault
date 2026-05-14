package cli

import (
	"bytes"
	"strings"
	"testing"
)

func TestDoctorAllPassedAfterInit(t *testing.T) {
	r := tempRunner(t)

	var out bytes.Buffer
	if code := runCommand(r, []string{"init"}, &out); code != 0 {
		t.Fatalf("init failed: %s", out.String())
	}
	out.Reset()

	// lock so vault exists
	out.Reset()
	if code := runCommand(r, []string{"lock"}, &out); code != 0 {
		t.Fatalf("lock failed: %s", out.String())
	}
	out.Reset()

	code := runDoctor(r, nil, &out)
	output := out.String()

	if !strings.Contains(output, "✓ config file present") {
		t.Errorf("expected config check to pass, got:\n%s", output)
	}
	if !strings.Contains(output, "✓ key file present") {
		t.Errorf("expected key check to pass, got:\n%s", output)
	}
	if !strings.Contains(output, "✓ vault file present") {
		t.Errorf("expected vault check to pass, got:\n%s", output)
	}
	if code != 0 && !strings.Contains(output, "Some checks failed") {
		// clipboard / editor may not be available in CI — that is acceptable
		t.Logf("doctor returned non-zero (likely missing editor/clipboard in CI): %s", output)
	}
}

func TestDoctorFailsWithoutInit(t *testing.T) {
	r := tempRunner(t)

	var out bytes.Buffer
	code := runDoctor(r, nil, &out)
	output := out.String()

	if code != 1 {
		t.Errorf("expected exit code 1 before init, got %d", code)
	}
	if !strings.Contains(output, "✗ config file present") {
		t.Errorf("expected config check to fail, got:\n%s", output)
	}
	if !strings.Contains(output, "✗ key file present") {
		t.Errorf("expected key check to fail, got:\n%s", output)
	}
	if !strings.Contains(output, "run 'envault init'") {
		t.Errorf("expected hint about init, got:\n%s", output)
	}
}

func TestDoctorVaultMissingAfterInit(t *testing.T) {
	r := tempRunner(t)

	var out bytes.Buffer
	if code := runCommand(r, []string{"init"}, &out); code != 0 {
		t.Fatalf("init failed: %s", out.String())
	}
	out.Reset()

	code := runDoctor(r, nil, &out)
	output := out.String()

	if !strings.Contains(output, "✗ vault file present") {
		t.Errorf("expected vault check to fail before lock, got:\n%s", output)
	}
	if !strings.Contains(output, "run 'envault lock'") {
		t.Errorf("expected hint about lock, got:\n%s", output)
	}
	if code != 1 {
		t.Errorf("expected exit code 1, got %d", code)
	}
}

func TestDoctorOutputContainsSummaryLine(t *testing.T) {
	r := tempRunner(t)

	var out bytes.Buffer
	runDoctor(r, nil, &out)
	output := out.String()

	hasSummary := strings.Contains(output, "All checks passed.") ||
		strings.Contains(output, "Some checks failed.")
	if !hasSummary {
		t.Errorf("expected a summary line, got:\n%s", output)
	}
}
