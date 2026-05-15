package cli_test

import (
	"strings"
	"testing"
)

func TestCloneCopiesKey(t *testing.T) {
	r := tempRunner(t)

	if err := r.Run([]string{"envault", "init"}); err != nil {
		t.Fatalf("init: %v", err)
	}
	if err := r.Run([]string{"envault", "lock"}); err != nil {
		t.Fatalf("lock: %v", err)
	}
	if err := r.Run([]string{"envault", "add", "SECRET=hunter2"}); err != nil {
		t.Fatalf("add: %v", err)
	}

	if err := r.Run([]string{"envault", "clone", "SECRET", "SECRET_BACKUP"}); err != nil {
		t.Fatalf("clone: %v", err)
	}

	out := r.Output()
	if !strings.Contains(out, "SECRET → SECRET_BACKUP") {
		t.Errorf("expected clone confirmation, got: %s", out)
	}
}

func TestCloneFailsWhenSourceMissing(t *testing.T) {
	r := tempRunner(t)

	if err := r.Run([]string{"envault", "init"}); err != nil {
		t.Fatalf("init: %v", err)
	}
	if err := r.Run([]string{"envault", "lock"}); err != nil {
		t.Fatalf("lock: %v", err)
	}

	err := r.Run([]string{"envault", "clone", "MISSING", "DEST"})
	if err == nil {
		t.Fatal("expected error for missing source key")
	}
	if !strings.Contains(err.Error(), "not found") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestCloneFailsWhenDestinationExists(t *testing.T) {
	r := tempRunner(t)

	if err := r.Run([]string{"envault", "init"}); err != nil {
		t.Fatalf("init: %v", err)
	}
	if err := r.Run([]string{"envault", "lock"}); err != nil {
		t.Fatalf("lock: %v", err)
	}
	if err := r.Run([]string{"envault", "add", "A=foo"}); err != nil {
		t.Fatalf("add A: %v", err)
	}
	if err := r.Run([]string{"envault", "add", "B=bar"}); err != nil {
		t.Fatalf("add B: %v", err)
	}

	err := r.Run([]string{"envault", "clone", "A", "B"})
	if err == nil {
		t.Fatal("expected error when destination exists")
	}
	if !strings.Contains(err.Error(), "already exists") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestCloneForceOverwritesDestination(t *testing.T) {
	r := tempRunner(t)

	if err := r.Run([]string{"envault", "init"}); err != nil {
		t.Fatalf("init: %v", err)
	}
	if err := r.Run([]string{"envault", "lock"}); err != nil {
		t.Fatalf("lock: %v", err)
	}
	if err := r.Run([]string{"envault", "add", "A=newval"}); err != nil {
		t.Fatalf("add A: %v", err)
	}
	if err := r.Run([]string{"envault", "add", "B=oldval"}); err != nil {
		t.Fatalf("add B: %v", err)
	}

	if err := r.Run([]string{"envault", "clone", "--force", "A", "B"}); err != nil {
		t.Fatalf("clone --force: %v", err)
	}

	out := r.Output()
	if !strings.Contains(out, "A → B") {
		t.Errorf("expected clone confirmation, got: %s", out)
	}
}

func TestCloneFailsWithoutInit(t *testing.T) {
	r := tempRunner(t)

	err := r.Run([]string{"envault", "clone", "A", "B"})
	if err == nil {
		t.Fatal("expected error without init")
	}
}

func TestCloneSameKeyReturnsError(t *testing.T) {
	r := tempRunner(t)

	if err := r.Run([]string{"envault", "init"}); err != nil {
		t.Fatalf("init: %v", err)
	}
	if err := r.Run([]string{"envault", "lock"}); err != nil {
		t.Fatalf("lock: %v", err)
	}
	if err := r.Run([]string{"envault", "add", "KEY=val"}); err != nil {
		t.Fatalf("add: %v", err)
	}

	err := r.Run([]string{"envault", "clone", "KEY", "KEY"})
	if err == nil {
		t.Fatal("expected error when source == destination")
	}
	if !strings.Contains(err.Error(), "must differ") {
		t.Errorf("unexpected error: %v", err)
	}
}
