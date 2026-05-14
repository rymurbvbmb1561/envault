package cli

import (
	"strings"
	"testing"
)

func TestSnapshotListEmptyWhenNoSnapshots(t *testing.T) {
	r := tempRunner(t)

	r.Run([]string{"init"})
	r.Run([]string{"add", "K=v"})
	r.Run([]string{"lock"})

	r.Stdout.Reset()
	if code := r.Run([]string{"snapshot-list"}); code != 0 {
		t.Fatalf("snapshot-list failed: %s", r.Stderr)
	}

	out := r.Stdout.String()
	if !strings.Contains(out, "no snapshots found") {
		t.Errorf("expected 'no snapshots found', got: %s", out)
	}
}

func TestSnapshotRestoreFailsWithoutInit(t *testing.T) {
	r := tempRunner(t)
	if code := r.Run([]string{"snapshot-restore", "20240101T000000Z"}); code != 1 {
		t.Errorf("expected exit 1 without init")
	}
}

func TestSnapshotRestoreFailsWithMissingSnapshot(t *testing.T) {
	r := tempRunner(t)

	r.Run([]string{"init"})
	r.Run([]string{"add", "K=v"})
	r.Run([]string{"lock"})

	if code := r.Run([]string{"snapshot-restore", "nonexistent"}); code != 1 {
		t.Errorf("expected exit 1 for missing snapshot")
	}

	errOut := r.Stderr.String()
	if !strings.Contains(errOut, "snapshot not found") {
		t.Errorf("expected 'snapshot not found' error, got: %s", errOut)
	}
}

func TestSnapshotRestoreRequiresArg(t *testing.T) {
	r := tempRunner(t)

	r.Run([]string{"init"})

	if code := r.Run([]string{"snapshot-restore"}); code != 1 {
		t.Errorf("expected exit 1 when no arg provided")
	}

	errOut := r.Stderr.String()
	if !strings.Contains(errOut, "usage:") {
		t.Errorf("expected usage message, got: %s", errOut)
	}
}
