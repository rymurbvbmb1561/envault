package cli

import (
	"bytes"
	"strings"
	"testing"
)

func TestFmtSortsKeys(t *testing.T) {
	r := tempRunner(t)

	if err := r.Init(); err != nil {
		t.Fatalf("init: %v", err)
	}

	// Write keys out of order.
	raw := map[string]string{
		"ZEBRA":   "z",
		"ALPHA":   "a",
		"MIDDLE":  "m",
	}
	if err := r.Vault.Write(r.KeyStore, raw); err != nil {
		t.Fatalf("write: %v", err)
	}

	var buf bytes.Buffer
	app := newTestApp(r, &buf)
	if err := app.Run([]string{"envault", "fmt"}); err != nil {
		t.Fatalf("fmt: %v", err)
	}

	if !strings.Contains(buf.String(), "formatted") {
		t.Errorf("expected 'formatted' in output, got: %q", buf.String())
	}
}

func TestFmtCheckPassesWhenAlreadySorted(t *testing.T) {
	r := tempRunner(t)

	if err := r.Init(); err != nil {
		t.Fatalf("init: %v", err)
	}

	raw := map[string]string{"ALPHA": "a", "BETA": "b"}
	if err := r.Vault.Write(r.KeyStore, raw); err != nil {
		t.Fatalf("write: %v", err)
	}

	// Run fmt once to ensure sorted.
	var buf bytes.Buffer
	app := newTestApp(r, &buf)
	_ = app.Run([]string{"envault", "fmt"})

	// Now --check should succeed.
	buf.Reset()
	app2 := newTestApp(r, &buf)
	if err := app2.Run([]string{"envault", "fmt", "--check"}); err != nil {
		t.Fatalf("fmt --check on sorted vault: %v", err)
	}
	if !strings.Contains(buf.String(), "already formatted") {
		t.Errorf("expected 'already formatted', got: %q", buf.String())
	}
}

func TestFmtFailsWithoutInit(t *testing.T) {
	r := tempRunner(t)
	var buf bytes.Buffer
	app := newTestApp(r, &buf)
	if err := app.Run([]string{"envault", "fmt"}); err == nil {
		t.Error("expected error when not initialised")
	}
}

func TestFmtFailsWithoutVault(t *testing.T) {
	r := tempRunner(t)
	if err := r.Init(); err != nil {
		t.Fatalf("init: %v", err)
	}
	// Do not write any vault file.
	var buf bytes.Buffer
	app := newTestApp(r, &buf)
	if err := app.Run([]string{"envault", "fmt"}); err == nil {
		t.Error("expected error when vault missing")
	}
}
