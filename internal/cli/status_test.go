package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/user/envault/internal/config"
	"github.com/user/envault/internal/keystore"
	"github.com/user/envault/internal/vault"
)

func tempStatusRunner(t *testing.T) (*Runner, string) {
	t.Helper()
	dir := t.TempDir()
	cfg := config.DefaultConfig(dir)
	ks := keystore.New(cfg.KeyPath)
	v := vault.New(cfg.VaultPath)
	r := &Runner{keys: ks, vault: v, cfg: cfg, out: os.Stdout}
	return r, dir
}

func TestStatusNotInitialised(t *testing.T) {
	r, _ := tempStatusRunner(t)
	var buf bytes.Buffer
	r.SetOutput(&buf)

	if err := r.Status(); err != nil {
		t.Fatalf("Status() unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "not initialised") {
		t.Errorf("expected 'not initialised' in output, got: %s", out)
	}
}

func TestStatusAfterInit(t *testing.T) {
	r, dir := tempStatusRunner(t)
	var buf bytes.Buffer
	r.SetOutput(&buf)

	// Simulate init by generating a key.
	if err := r.keys.Generate(); err != nil {
		t.Fatalf("Generate() error: %v", err)
	}

	if err := r.Status(); err != nil {
		t.Fatalf("Status() unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "found") {
		t.Errorf("expected 'found' key status, got: %s", out)
	}
	_ = dir
}

func TestVersionOutput(t *testing.T) {
	r, _ := tempStatusRunner(t)
	var buf bytes.Buffer
	r.SetOutput(&buf)

	if err := r.Version("1.2.3"); err != nil {
		t.Fatalf("Version() unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "1.2.3") {
		t.Errorf("expected version in output, got: %s", buf.String())
	}
}

func TestVersionDefaultsDev(t *testing.T) {
	r, _ := tempStatusRunner(t)
	var buf bytes.Buffer
	r.SetOutput(&buf)

	_ = r.Version("")
	if !strings.Contains(buf.String(), "dev") {
		t.Errorf("expected 'dev' version, got: %s", buf.String())
	}
}

func TestRequireInitReturnsError(t *testing.T) {
	r, _ := tempStatusRunner(t)
	if err := r.requireInit(); err == nil {
		t.Error("expected error when not initialised")
	}
}

func TestRequireVaultReturnsErrorWithoutVault(t *testing.T) {
	r, dir := tempStatusRunner(t)
	_ = r.keys.Generate()
	if err := r.requireVault(); err == nil {
		t.Error("expected error when vault is absent")
	}
	_ = filepath.Join(dir, ".envault")
	_ = os.RemoveAll(dir)
}
