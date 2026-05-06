package cli_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/nicholasgasior/envault/internal/cli"
)

// writeEditorScript writes a small shell script that appends a line to the
// file it receives, simulating a user edit, and returns the script path.
func writeEditorScript(t *testing.T, dir, extraLine string) string {
	t.Helper()
	script := filepath.Join(dir, "fake-editor.sh")
	body := "#!/bin/sh\necho '" + extraLine + "' >> \"$1\"\n"
	if err := os.WriteFile(script, []byte(body), 0o755); err != nil {
		t.Fatalf("writeEditorScript: %v", err)
	}
	return script
}

func TestEditRoundTrip(t *testing.T) {
	r, dir := tempRunner(t)

	if err := r.Init(); err != nil {
		t.Fatalf("Init: %v", err)
	}

	// Create an initial .env and lock it.
	envPath := filepath.Join(dir, ".env")
	if err := os.WriteFile(envPath, []byte("FOO=bar\n"), 0o600); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}
	if err := r.Lock(envPath); err != nil {
		t.Fatalf("Lock: %v", err)
	}

	// Point EDITOR at a script that appends a line.
	script := writeEditorScript(t, dir, "EDITED=1")
	t.Setenv("EDITOR", script)

	if err := r.Edit(); err != nil {
		t.Fatalf("Edit: %v", err)
	}

	// Unlock and verify the new line is present.
	outPath := filepath.Join(dir, ".env.out")
	if err := r.Unlock(outPath); err != nil {
		t.Fatalf("Unlock: %v", err)
	}

	data, err := os.ReadFile(outPath)
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}
	if !strings.Contains(string(data), "EDITED=1") {
		t.Errorf("expected EDITED=1 in decrypted output, got: %s", data)
	}
	if !strings.Contains(string(data), "FOO=bar") {
		t.Errorf("expected FOO=bar preserved in output, got: %s", data)
	}
}

func TestEditFailsWithoutVault(t *testing.T) {
	r, _ := tempRunner(t)

	if err := r.Init(); err != nil {
		t.Fatalf("Init: %v", err)
	}

	err := r.Edit()
	if err == nil {
		t.Fatal("expected error when vault does not exist, got nil")
	}
}

func TestEditFailsWithoutInit(t *testing.T) {
	r, _ := tempRunner(t)

	err := r.Edit()
	if err == nil {
		t.Fatal("expected error when not initialised, got nil")
	}
	if !strings.Contains(err.Error(), "not initialised") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestEditFailsWithNoEditor(t *testing.T) {
	r, dir := tempRunner(t)

	if err := r.Init(); err != nil {
		t.Fatalf("Init: %v", err)
	}

	envPath := filepath.Join(dir, ".env")
	if err := os.WriteFile(envPath, []byte("X=1\n"), 0o600); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}
	if err := r.Lock(envPath); err != nil {
		t.Fatalf("Lock: %v", err)
	}

	// Unset EDITOR and ensure no default editor is found by using a PATH that
	// contains no executables.
	t.Setenv("EDITOR", "")
	t.Setenv("PATH", t.TempDir())

	err := r.Edit()
	if err == nil {
		t.Fatal("expected error when no editor is available, got nil")
	}
}

// Compile-time check: cli.Runner must expose Edit.
var _ = (*cli.Runner)(nil)
