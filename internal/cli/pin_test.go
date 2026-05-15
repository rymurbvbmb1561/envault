package cli

import (
	"strings"
	"testing"
)

func TestPinAddAndList(t *testing.T) {
	r := tempRunner(t)
	if err := r.Run([]string{"envault", "init"}); err != nil {
		t.Fatalf("init: %v", err)
	}
	if err := r.Run([]string{"envault", "add", "API_KEY=secret"}); err != nil {
		t.Fatalf("add: %v", err)
	}
	if err := r.Run([]string{"envault", "pin", "add", "API_KEY"}); err != nil {
		t.Fatalf("pin add: %v", err)
	}
	out := r.Output()
	if !strings.Contains(out, "pinned API_KEY") {
		t.Errorf("expected pinned confirmation, got: %s", out)
	}

	if err := r.Run([]string{"envault", "pin", "list"}); err != nil {
		t.Fatalf("pin list: %v", err)
	}
	out = r.Output()
	if !strings.Contains(out, "API_KEY") {
		t.Errorf("expected API_KEY in pin list, got: %s", out)
	}
}

func TestPinRemove(t *testing.T) {
	r := tempRunner(t)
	if err := r.Run([]string{"envault", "init"}); err != nil {
		t.Fatalf("init: %v", err)
	}
	if err := r.Run([]string{"envault", "add", "DB_URL=postgres://localhost"}); err != nil {
		t.Fatalf("add: %v", err)
	}
	if err := r.Run([]string{"envault", "pin", "add", "DB_URL"}); err != nil {
		t.Fatalf("pin add: %v", err)
	}
	if err := r.Run([]string{"envault", "pin", "remove", "DB_URL"}); err != nil {
		t.Fatalf("pin remove: %v", err)
	}
	out := r.Output()
	if !strings.Contains(out, "unpinned DB_URL") {
		t.Errorf("expected unpin confirmation, got: %s", out)
	}

	if err := r.Run([]string{"envault", "pin", "list"}); err != nil {
		t.Fatalf("pin list: %v", err)
	}
	out = r.Output()
	if !strings.Contains(out, "no pinned keys") {
		t.Errorf("expected no pinned keys, got: %s", out)
	}
}

func TestPinListEmptyWhenNoPins(t *testing.T) {
	r := tempRunner(t)
	if err := r.Run([]string{"envault", "init"}); err != nil {
		t.Fatalf("init: %v", err)
	}
	if err := r.Run([]string{"envault", "pin", "list"}); err != nil {
		t.Fatalf("pin list: %v", err)
	}
	out := r.Output()
	if !strings.Contains(out, "no pinned keys") {
		t.Errorf("expected no pinned keys message, got: %s", out)
	}
}

func TestPinAddMissingKeyReturnsError(t *testing.T) {
	r := tempRunner(t)
	if err := r.Run([]string{"envault", "init"}); err != nil {
		t.Fatalf("init: %v", err)
	}
	err := r.Run([]string{"envault", "pin", "add", "NONEXISTENT"})
	if err == nil {
		t.Fatal("expected error for missing key, got nil")
	}
	if !strings.Contains(err.Error(), "not found") {
		t.Errorf("expected 'not found' error, got: %v", err)
	}
}

func TestPinFailsWithoutInit(t *testing.T) {
	r := tempRunner(t)
	err := r.Run([]string{"envault", "pin", "add", "KEY"})
	if err == nil {
		t.Fatal("expected error without init, got nil")
	}
}

func TestPinAddNoArgsReturnsError(t *testing.T) {
	r := tempRunner(t)
	if err := r.Run([]string{"envault", "init"}); err != nil {
		t.Fatalf("init: %v", err)
	}
	err := r.Run([]string{"envault", "pin", "add"})
	if err == nil {
		t.Fatal("expected error with no args, got nil")
	}
}
