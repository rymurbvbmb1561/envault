package cli

import (
	"bytes"
	"strings"
	"testing"
)

func TestCompareNoChanges(t *testing.T) {
	r := tempRunner(t)
	runInit(r.cliCtx())
	runLock(r.cliCtx())

	// create a snapshot then compare it against the current vault (identical)
	runSnapshot(r.cliCtxWithArgs("before"))

	var buf bytes.Buffer
	r.App.Writer = &buf
	err := runCompare(r.cliCtxWithArgs("before"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "no differences found") {
		t.Errorf("expected 'no differences found', got: %s", buf.String())
	}
}

func TestCompareDetectsAddition(t *testing.T) {
	r := tempRunner(t)
	runInit(r.cliCtx())
	runLock(r.cliCtx())
	runSnapshot(r.cliCtxWithArgs("before"))

	// add a key after snapshot
	runAdd(r.cliCtxWithArgs("NEW_KEY=newvalue"))
	runLock(r.cliCtx())

	var buf bytes.Buffer
	r.App.Writer = &buf
	err := runCompare(r.cliCtxWithArgs("before"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "+ NEW_KEY") {
		t.Errorf("expected addition marker, got: %s", out)
	}
}

func TestCompareDetectsRemoval(t *testing.T) {
	r := tempRunner(t)
	runInit(r.cliCtx())
	runAdd(r.cliCtxWithArgs("OLD_KEY=oldvalue"))
	runLock(r.cliCtx())
	runSnapshot(r.cliCtxWithArgs("before"))

	runRemove(r.cliCtxWithArgs("OLD_KEY"))
	runLock(r.cliCtx())

	var buf bytes.Buffer
	r.App.Writer = &buf
	err := runCompare(r.cliCtxWithArgs("before"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "- OLD_KEY") {
		t.Errorf("expected removal marker, got: %s", out)
	}
}

func TestCompareFailsWithoutInit(t *testing.T) {
	r := tempRunner(t)
	err := runCompare(r.cliCtxWithArgs("snap"))
	if err == nil {
		t.Fatal("expected error when not initialised")
	}
}

func TestCompareRequiresArg(t *testing.T) {
	r := tempRunner(t)
	runInit(r.cliCtx())
	err := runCompare(r.cliCtx())
	if err == nil {
		t.Fatal("expected error when no snapshot arg provided")
	}
	if !strings.Contains(err.Error(), "usage") {
		t.Errorf("expected usage hint in error, got: %v", err)
	}
}

func TestCompareTwoSnapshots(t *testing.T) {
	r := tempRunner(t)
	runInit(r.cliCtx())
	runLock(r.cliCtx())
	runSnapshot(r.cliCtxWithArgs("snap-a"))

	runAdd(r.cliCtxWithArgs("EXTRA=1"))
	runLock(r.cliCtx())
	runSnapshot(r.cliCtxWithArgs("snap-b"))

	var buf bytes.Buffer
	r.App.Writer = &buf
	err := runCompare(r.cliCtxWithArgs("snap-a", "snap-b"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "+ EXTRA") {
		t.Errorf("expected EXTRA to appear as addition, got: %s", out)
	}
}
