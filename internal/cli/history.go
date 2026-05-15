package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/urfave/cli/v2"
)

// HistoryEntry records a single mutation event for a secret key.
type HistoryEntry struct {
	Timestamp time.Time `json:"timestamp"`
	Action    string    `json:"action"` // set | remove | rename
	Key       string    `json:"key"`
	PrevKey   string    `json:"prev_key,omitempty"`
}

// AppendHistory writes a new entry to the history log file alongside the vault.
func AppendHistory(vaultPath string, entry HistoryEntry) error {
	path := historyPath(vaultPath)

	var entries []HistoryEntry
	if data, err := os.ReadFile(path); err == nil {
		_ = json.Unmarshal(data, &entries)
	}

	entries = append(entries, entry)

	data, err := json.MarshalIndent(entries, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal history: %w", err)
	}
	return os.WriteFile(path, data, 0600)
}

// LoadHistory reads all history entries for the given vault.
func LoadHistory(vaultPath string) ([]HistoryEntry, error) {
	path := historyPath(vaultPath)
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("read history: %w", err)
	}
	var entries []HistoryEntry
	if err := json.Unmarshal(data, &entries); err != nil {
		return nil, fmt.Errorf("parse history: %w", err)
	}
	return entries, nil
}

func historyPath(vaultPath string) string {
	dir := filepath.Dir(vaultPath)
	return filepath.Join(dir, ".envault_history.json")
}

func runHistory(c *cli.Context) error {
	r, err := newRunnerFromContext(c)
	if err != nil {
		return err
	}
	if err := requireInitialised(r); err != nil {
		return err
	}

	entries, err := LoadHistory(r.VaultPath())
	if err != nil {
		return fmt.Errorf("load history: %w", err)
	}

	filterKey := c.String("key")
	limit := c.Int("limit")

	printed := 0
	for i := len(entries) - 1; i >= 0 && printed < limit; i-- {
		e := entries[i]
		if filterKey != "" && e.Key != filterKey && e.PrevKey != filterKey {
			continue
		}
		ts := e.Timestamp.Format("2006-01-02 15:04:05")
		switch e.Action {
		case "rename":
			fmt.Fprintf(c.App.Writer, "%s  %-8s  %s -> %s\n", ts, e.Action, e.PrevKey, e.Key)
		default:
			fmt.Fprintf(c.App.Writer, "%s  %-8s  %s\n", ts, e.Action, e.Key)
		}
		printed++
	}

	if printed == 0 {
		fmt.Fprintln(c.App.Writer, "No history found.")
	}
	return nil
}
