package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/urfave/cli/v2"
)

func TestWatchWritesOutputFileOnStart(t *testing.T) {
	r := tempRunner(t)
	if err := r.Init(); err != nil {
		t.Fatalf("init: %v", err)
	}
	if err := r.Lock(); err != nil {
		t.Fatalf("lock: %v", err)
	}

	// Pre-populate a secret.
	data := map[string]string{"WATCH_KEY": "hello"}
	if err := r.Vault.Write(r.Key, data); err != nil {
		t.Fatalf("write vault: %v", err)
	}

	outFile := filepath.Join(t.TempDir(), "out.env")

	var buf bytes.Buffer
	app := cli.NewApp()
	app.Writer = &buf
	app.ErrWriter = &buf
	app.Commands = []*cli.Command{watchCommand()}

	// Run watch in a goroutine and cancel via the output file being written.
	done := make(chan error, 1)
	go func() {
		done <- app.Run([]string{"envault", "watch", "--output", outFile})
	}()

	// Give the watcher time to do the initial write.
	time.Sleep(200 * time.Millisecond)

	content, err := os.ReadFile(outFile)
	if err != nil {
		t.Fatalf("output file not created: %v", err)
	}
	if !strings.Contains(string(content), "WATCH_KEY=hello") {
		t.Errorf("expected WATCH_KEY=hello in output, got: %s", content)
	}
}

func TestWatchFailsWithoutInit(t *testing.T) {
	r := tempRunner(t)
	_ = r

	var buf bytes.Buffer
	app := cli.NewApp()
	app.Writer = &buf
	app.ErrWriter = &buf
	app.Commands = []*cli.Command{watchCommand()}

	err := app.Run([]string{"envault", "watch"})
	if err == nil {
		t.Fatal("expected error when not initialised")
	}
}

func TestWatchFailsWithoutVault(t *testing.T) {
	r := tempRunner(t)
	if err := r.Init(); err != nil {
		t.Fatalf("init: %v", err)
	}
	// Do NOT call r.Lock() so no vault exists.

	var buf bytes.Buffer
	app := cli.NewApp()
	app.Writer = &buf
	app.ErrWriter = &buf
	app.Commands = []*cli.Command{watchCommand()}

	err := app.Run([]string{"envault", "watch"})
	if err == nil {
		t.Fatal("expected error when vault missing")
	}
}
