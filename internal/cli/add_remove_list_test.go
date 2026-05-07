package cli

import (
	"bytes"
	"strings"
	"testing"
)

func TestAddAndListRoundTrip(t *testing.T) {
	r := tempRunner(t)

	if err := r.initCmd(nil); err != nil {
		t.Fatalf("init: %v", err)
	}

	if err := r.addCmd([]string{"FOO=bar"}); err != nil {
		t.Fatalf("add: %v", err)
	}

	var buf bytes.Buffer
	r.stdout = &buf
	if err := r.listCmd(nil); err != nil {
		t.Fatalf("list: %v", err)
	}

	if !strings.Contains(buf.String(), "FOO") {
		t.Errorf("expected FOO in list output, got: %s", buf.String())
	}
}

func TestAddUpdatesExistingKey(t *testing.T) {
	r := tempRunner(t)
	if err := r.initCmd(nil); err != nil {
		t.Fatalf("init: %v", err)
	}

	_ = r.addCmd([]string{"KEY=first"})
	_ = r.addCmd([]string{"KEY=second"})

	recipient, _ := r.keystore.Load()
	env, err := r.vault.ReadEnv(recipient)
	if err != nil {
		t.Fatalf("read: %v", err)
	}
	if env["KEY"] != "second" {
		t.Errorf("expected 'second', got %q", env["KEY"])
	}
}

func TestRemoveDeletesKey(t *testing.T) {
	r := tempRunner(t)
	if err := r.initCmd(nil); err != nil {
		t.Fatalf("init: %v", err)
	}

	_ = r.addCmd([]string{"SECRET=value"})

	if err := r.removeCmd([]string{"SECRET"}); err != nil {
		t.Fatalf("remove: %v", err)
	}

	recipient, _ := r.keystore.Load()
	env, err := r.vault.ReadEnv(recipient)
	if err != nil {
		t.Fatalf("read: %v", err)
	}
	if _, ok := env["SECRET"]; ok {
		t.Error("expected SECRET to be removed")
	}
}

func TestRemoveMissingKeyReturnsError(t *testing.T) {
	r := tempRunner(t)
	if err := r.initCmd(nil); err != nil {
		t.Fatalf("init: %v", err)
	}
	_ = r.addCmd([]string{"A=1"})

	if err := r.removeCmd([]string{"MISSING"}); err == nil {
		t.Error("expected error for missing key")
	}
}

func TestAddInvalidFormatReturnsError(t *testing.T) {
	r := tempRunner(t)
	if err := r.initCmd(nil); err != nil {
		t.Fatalf("init: %v", err)
	}

	if err := r.addCmd([]string{"NOEQUALS"}); err == nil {
		t.Error("expected error for missing '='")
	}
}

func TestListEmptyVaultMessage(t *testing.T) {
	r := tempRunner(t)
	if err := r.initCmd(nil); err != nil {
		t.Fatalf("init: %v", err)
	}

	var buf bytes.Buffer
	r.stdout = &buf
	if err := r.listCmd(nil); err != nil {
		t.Fatalf("list: %v", err)
	}

	if !strings.Contains(buf.String(), "no secrets") {
		t.Errorf("expected empty message, got: %s", buf.String())
	}
}
