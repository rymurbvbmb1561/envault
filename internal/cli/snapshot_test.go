package cli

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestSnapshotCreatesFile(t *testing.T) {
	r := tempRunner(t)

	if code := r.Run([]string{"init"}); code != 0 {
		t.Fatalf("init failed: %s", r.Stderr)
	}
	if code := r.Run([]string{"add", "KEY=value"}); code != 0 {
		t.Fatalf("add failed: %s", r.Stderr)
	}
	if code := r.Run([]string{"lock"}); code != 0 {
		t.Fatalf("lock failed: %s", r.Stderr)
	}

	if code := r.Run([]string{"snapshot"}); code != 0 {
		t.Fatalf("snapshot failed: %s", r.Stderr)
	}

	out := r.Stdout.String()
	if !strings.Contains(out, "snapshot saved:") {
		t.Errorf("expected 'snapshot saved:' in output, got: %s", out)
	}
}

func TestSnapshotFailsWithoutVault(t *testing.T) {
	r := tempRunner(t)

	if code := r.Run([]string{"init"}); code != 0 {
		t.Fatalf("init failed")
	}

	if code := r.Run([]string{"snapshot"}); code != 1 {
		t.Errorf("expected exit 1 without vault, got 0")
	}
}

func TestSnapshotFailsWithoutInit(t *testing.T) {
	r := tempRunner(t)
	if code := r.Run([]string{"snapshot"}); code != 1 {
		t.Errorf("expected exit 1 without init")
	}
}

func TestSnapshotListShowsSnapshots(t *testing.T) {
	r := tempRunner(t)

	r.Run([]string{"init"})
	r.Run([]string{"add", "X=1"})
	r.Run([]string{"lock"})
	r.Run([]string{"snapshot"})

	r.Stdout.Reset()
	if code := r.Run([]string{"snapshot-list"}); code != 0 {
		t.Fatalf("snapshot-list failed: %s", r.Stderr)
	}

	out := r.Stdout.String()
	if !strings.Contains(out, ".env") {
		t.Errorf("expected .env file in list, got: %s", out)
	}
}

func TestSnapshotRestoreRoundTrip(t *testing.T) {
	r := tempRunner(t)

	r.Run([]string{"init"})
	r.Run([]string{"add", "RESTORE_KEY=hello"})
	r.Run([]string{"lock"})
	r.Run([]string{"snapshot"})

	// find the snapshot file name
	cfg := r.cfg
	snapshotDir := filepath.Join(filepath.Dir(cfg.VaultPath), "snapshots")
	entries, err := os.ReadDir(snapshotDir)
	if err != nil || len(entries) == 0 {
		t.Fatalf("no snapshots found")
	}
	name := strings.TrimSuffix(entries[0].Name(), ".env")

	// overwrite vault with different content
	r.Run([]string{"add", "RESTORE_KEY=changed"})
	r.Run([]string{"lock"})

	// restore
	if code := r.Run([]string{"snapshot-restore", name}); code != 0 {
		t.Fatalf("snapshot-restore failed: %s", r.Stderr)
	}

	r.Stdout.Reset()
	r.Run([]string{"unlock"})

	// The restored .env should contain the original value
	envPath := filepath.Join(filepath.Dir(cfg.VaultPath), ".env")
	content, _ := os.ReadFile(envPath)
	if !strings.Contains(string(content), "RESTORE_KEY=hello") {
		t.Errorf("expected restored value 'hello', got: %s", string(content))
	}
}
