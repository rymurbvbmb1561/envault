package cli_test

import (
	"strings"
	"testing"
)

func TestCloneRequiresTwoArgs(t *testing.T) {
	r := tempRunner(t)

	if err := r.Run([]string{"envault", "init"}); err != nil {
		t.Fatalf("init: %v", err)
	}
	if err := r.Run([]string{"envault", "lock"}); err != nil {
		t.Fatalf("lock: %v", err)
	}

	err := r.Run([]string{"envault", "clone", "ONLY_ONE"})
	if err == nil {
		t.Fatal("expected error with only one arg")
	}
	if !strings.Contains(err.Error(), "usage") {
		t.Errorf("expected usage hint, got: %v", err)
	}
}

func TestCloneFailsWithoutVault(t *testing.T) {
	r := tempRunner(t)

	if err := r.Run([]string{"envault", "init"}); err != nil {
		t.Fatalf("init: %v", err)
	}
	// deliberately skip 'lock' so no vault file exists

	err := r.Run([]string{"envault", "clone", "A", "B"})
	if err == nil {
		t.Fatal("expected error without vault")
	}
	if !strings.Contains(err.Error(), "no vault found") {
		t.Errorf("unexpected error: %v", err)
	}
}
