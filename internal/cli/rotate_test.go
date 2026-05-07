package cli

import (
	"bytes"
	"strings"
	"testing"
)

func TestRotateReEncryptsVault(t *testing.T) {
	r := tempRunner(t)

	// Initialise and lock some secrets.
	if err := r.init(&bytes.Buffer{}); err != nil {
		t.Fatalf("init: %v", err)
	}
	if err := r.lock(&bytes.Buffer{}); err != nil {
		t.Fatalf("lock: %v", err)
	}

	// Capture the old recipient string before rotation.
	oldIdentity, err := r.keystore.Load()
	if err != nil {
		t.Fatalf("load old identity: %v", err)
	}
	oldRecipient := oldIdentity.Recipient().String()

	var out bytes.Buffer
	if err := r.rotate(&out); err != nil {
		t.Fatalf("rotate: %v", err)
	}

	if !strings.Contains(out.String(), "rotated") {
		t.Errorf("expected success message, got: %q", out.String())
	}

	// The key on disk must have changed.
	newIdentity, err := r.keystore.Load()
	if err != nil {
		t.Fatalf("load new identity: %v", err)
	}
	if newIdentity.Recipient().String() == oldRecipient {
		t.Error("expected a new recipient after rotation, but it is unchanged")
	}

	// The vault must still be readable with the new key.
	plaintext, err := r.vault.Read(newIdentity)
	if err != nil {
		t.Fatalf("reading vault after rotation: %v", err)
	}
	if len(plaintext) == 0 {
		t.Error("expected non-empty plaintext after rotation")
	}
}

func TestRotateFailsWithoutVault(t *testing.T) {
	r := tempRunner(t)
	if err := r.init(&bytes.Buffer{}); err != nil {
		t.Fatalf("init: %v", err)
	}
	// No lock step – vault does not exist.
	if err := r.rotate(&bytes.Buffer{}); err == nil {
		t.Error("expected error rotating without a vault")
	}
}

func TestRotateFailsWithoutInit(t *testing.T) {
	r := tempRunner(t)
	if err := r.rotate(&bytes.Buffer{}); err == nil {
		t.Error("expected error rotating without init")
	}
}
