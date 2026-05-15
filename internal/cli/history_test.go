package cli

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func TestAppendAndLoadHistory(t *testing.T) {
	r := tempRunner(t)
	vaultPath := r.VaultPath()

	entry := HistoryEntry{
		Timestamp: time.Now().UTC(),
		Action:    "set",
		Key:       "API_KEY",
	}
	if err := AppendHistory(vaultPath, entry); err != nil {
		t.Fatalf("AppendHistory: %v", err)
	}

	entries, err := LoadHistory(vaultPath)
	if err != nil {
		t.Fatalf("LoadHistory: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if entries[0].Key != "API_KEY" {
		t.Errorf("expected key API_KEY, got %s", entries[0].Key)
	}
}

func TestLoadHistoryEmptyWhenMissing(t *testing.T) {
	r := tempRunner(t)
	entries, err := LoadHistory(r.VaultPath())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 0 {
		t.Errorf("expected empty history, got %d entries", len(entries))
	}
}

func TestHistoryCommandPrintsEntries(t *testing.T) {
	r := tempRunner(t)
	if err := r.Init(); err != nil {
		t.Fatalf("init: %v", err)
	}

	now := time.Now().UTC()
	_ = AppendHistory(r.VaultPath(), HistoryEntry{Timestamp: now, Action: "set", Key: "DB_URL"})
	_ = AppendHistory(r.VaultPath(), HistoryEntry{Timestamp: now, Action: "remove", Key: "OLD_KEY"})

	var buf bytes.Buffer
	err := runApp(r, &buf, "history")
	if err != nil {
		t.Fatalf("history command: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "DB_URL") {
		t.Errorf("expected DB_URL in output, got:\n%s", out)
	}
	if !strings.Contains(out, "OLD_KEY") {
		t.Errorf("expected OLD_KEY in output, got:\n%s", out)
	}
}

func TestHistoryFilterByKey(t *testing.T) {
	r := tempRunner(t)
	if err := r.Init(); err != nil {
		t.Fatalf("init: %v", err)
	}

	now := time.Now().UTC()
	_ = AppendHistory(r.VaultPath(), HistoryEntry{Timestamp: now, Action: "set", Key: "DB_URL"})
	_ = AppendHistory(r.VaultPath(), HistoryEntry{Timestamp: now, Action: "set", Key: "SECRET_TOKEN"})

	var buf bytes.Buffer
	err := runApp(r, &buf, "history", "--key", "DB_URL")
	if err != nil {
		t.Fatalf("history command: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "DB_URL") {
		t.Errorf("expected DB_URL in output")
	}
	if strings.Contains(out, "SECRET_TOKEN") {
		t.Errorf("SECRET_TOKEN should be filtered out")
	}
}

func TestHistoryEmptyPrintsMessage(t *testing.T) {
	r := tempRunner(t)
	if err := r.Init(); err != nil {
		t.Fatalf("init: %v", err)
	}

	var buf bytes.Buffer
	if err := runApp(r, &buf, "history"); err != nil {
		t.Fatalf("history command: %v", err)
	}

	if !strings.Contains(buf.String(), "No history found") {
		t.Errorf("expected 'No history found', got: %s", buf.String())
	}
}
