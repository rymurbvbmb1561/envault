package cli_test

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
)

func TestExportPrintsExportPrefixByDefault(t *testing.T) {
	r, _, cleanup := tempRunner(t)
	defer cleanup()

	if err := r.Init(); err != nil {
		t.Fatalf("init: %v", err)
	}
	if err := r.Lock(); err != nil {
		t.Fatalf("lock: %v", err)
	}

	// Capture stdout.
	old := os.Stdout
	rp, wp, _ := os.Pipe()
	os.Stdout = wp

	err := r.Export("export")

	wp.Close()
	os.Stdout = old

	var buf bytes.Buffer
	io.Copy(&buf, rp)

	if err != nil {
		t.Fatalf("export: %v", err)
	}

	output := buf.String()
	for _, line := range strings.Split(strings.TrimSpace(output), "\n") {
		if line == "" {
			continue
		}
		if !strings.HasPrefix(line, "export ") {
			t.Errorf("expected 'export ' prefix, got: %q", line)
		}
	}
}

func TestExportDotenvFormat(t *testing.T) {
	r, _, cleanup := tempRunner(t)
	defer cleanup()

	if err := r.Init(); err != nil {
		t.Fatalf("init: %v", err)
	}
	if err := r.Lock(); err != nil {
		t.Fatalf("lock: %v", err)
	}

	old := os.Stdout
	rp, wp, _ := os.Pipe()
	os.Stdout = wp

	err := r.Export("dotenv")

	wp.Close()
	os.Stdout = old

	var buf bytes.Buffer
	io.Copy(&buf, rp)

	if err != nil {
		t.Fatalf("export dotenv: %v", err)
	}

	output := buf.String()
	for _, line := range strings.Split(strings.TrimSpace(output), "\n") {
		if line == "" {
			continue
		}
		if strings.HasPrefix(line, "export ") {
			t.Errorf("dotenv format should not have 'export ' prefix, got: %q", line)
		}
	}
}

func TestExportFailsWithoutVault(t *testing.T) {
	r, _, cleanup := tempRunner(t)
	defer cleanup()

	if err := r.Init(); err != nil {
		t.Fatalf("init: %v", err)
	}

	if err := r.Export("export"); err == nil {
		t.Fatal("expected error when vault does not exist")
	}
}

func TestExportFailsWithoutInit(t *testing.T) {
	r, _, cleanup := tempRunner(t)
	defer cleanup()

	if err := r.Export("export"); err == nil {
		t.Fatal("expected error when not initialised")
	}
}
