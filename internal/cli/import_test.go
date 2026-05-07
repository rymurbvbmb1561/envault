package cli

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestImportAddsKeysToVault(t *testing.T) {
	r := tempRunner(t)
	if code := r.Run([]string{"init"}); code != 0 {
		t.Fatal("init failed")
	}

	envFile := filepath.Join(t.TempDir(), ".env")
	content := "FOO=bar\nBAZ=qux\n# comment\n\nEMPTY=\n"
	if err := os.WriteFile(envFile, []byte(content), 0600); err != nil {
		t.Fatal(err)
	}

	if code := r.Run([]string{"import", envFile}); code != 0 {
		t.Fatalf("import failed: %s", r.Stderr)
	}

	out := r.Stdout.String()
	if !strings.Contains(out, "imported 3 key(s)") {
		t.Errorf("unexpected output: %q", out)
	}

	// Verify keys are readable via list
	r.Stdout.Reset()
	if code := r.Run([]string{"list"}); code != 0 {
		t.Fatal("list failed")
	}
	listOut := r.Stdout.String()
	for _, key := range []string{"FOO", "BAZ", "EMPTY"} {
		if !strings.Contains(listOut, key) {
			t.Errorf("expected key %q in list output, got: %q", key, listOut)
		}
	}
}

func TestImportMergesWithExistingKeys(t *testing.T) {
	r := tempRunner(t)
	if code := r.Run([]string{"init"}); code != 0 {
		t.Fatal("init failed")
	}
	r.Run([]string{"add", "EXISTING=original"})

	envFile := filepath.Join(t.TempDir(), ".env")
	os.WriteFile(envFile, []byte("EXISTING=updated\nNEW=value\n"), 0600)

	if code := r.Run([]string{"import", envFile}); code != 0 {
		t.Fatalf("import failed: %s", r.Stderr)
	}

	r.Stdout.Reset()
	r.Run([]string{"export", "--format", "dotenv"})
	out := r.Stdout.String()
	if !strings.Contains(out, "EXISTING=updated") {
		t.Errorf("expected EXISTING to be updated, got: %q", out)
	}
	if !strings.Contains(out, "NEW=value") {
		t.Errorf("expected NEW key, got: %q", out)
	}
}

func TestImportFailsWithoutInit(t *testing.T) {
	r := tempRunner(t)
	envFile := filepath.Join(t.TempDir(), ".env")
	os.WriteFile(envFile, []byte("K=V\n"), 0600)

	if code := r.Run([]string{"import", envFile}); code == 0 {
		t.Error("expected non-zero exit without init")
	}
}

func TestImportFailsWithNoArgs(t *testing.T) {
	r := tempRunner(t)
	r.Run([]string{"init"})

	if code := r.Run([]string{"import"}); code == 0 {
		t.Error("expected non-zero exit with no args")
	}
	if !strings.Contains(r.Stderr.String(), "usage") {
		t.Errorf("expected usage message, got: %q", r.Stderr.String())
	}
}

func TestImportFailsWithMissingFile(t *testing.T) {
	r := tempRunner(t)
	r.Run([]string{"init"})

	if code := r.Run([]string{"import", "/nonexistent/.env"}); code == 0 {
		t.Error("expected non-zero exit for missing file")
	}
}
