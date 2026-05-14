package cli

import (
	"bytes"
	"strings"
	"testing"
)

func setupGrepVault(t *testing.T) *runner {
	t.Helper()
	r := tempRunner(t)
	if code := r.Run([]string{"init"}); code != 0 {
		t.Fatal("init failed")
	}
	for _, kv := range []string{
		"DATABASE_URL=postgres://localhost/mydb",
		"REDIS_URL=redis://localhost:6379",
		"SECRET_KEY=supersecret",
		"API_TOKEN=abc123",
	} {
		if code := r.Run([]string{"add", kv}); code != 0 {
			t.Fatalf("add %q failed", kv)
		}
	}
	return r
}

func TestGrepMatchesValue(t *testing.T) {
	r := setupGrepVault(t)
	var buf bytes.Buffer
	r.Stdout = &buf
	code := r.Run([]string{"grep", "localhost"})
	if code != 0 {
		t.Fatalf("expected exit 0, got %d", code)
	}
	out := buf.String()
	if !strings.Contains(out, "DATABASE_URL") {
		t.Error("expected DATABASE_URL in output")
	}
	if !strings.Contains(out, "REDIS_URL") {
		t.Error("expected REDIS_URL in output")
	}
	if strings.Contains(out, "SECRET_KEY") {
		t.Error("unexpected SECRET_KEY in output")
	}
}

func TestGrepMatchesKeyWithFlag(t *testing.T) {
	r := setupGrepVault(t)
	var buf bytes.Buffer
	r.Stdout = &buf
	code := r.Run([]string{"grep", "-keys", "TOKEN"})
	if code != 0 {
		t.Fatalf("expected exit 0, got %d", code)
	}
	out := buf.String()
	if !strings.Contains(out, "API_TOKEN") {
		t.Error("expected API_TOKEN in output")
	}
}

func TestGrepCaseInsensitive(t *testing.T) {
	r := setupGrepVault(t)
	var buf bytes.Buffer
	r.Stdout = &buf
	code := r.Run([]string{"grep", "-i", "SUPERSECRET"})
	if code != 0 {
		t.Fatalf("expected exit 0, got %d", code)
	}
	out := buf.String()
	if !strings.Contains(out, "SECRET_KEY") {
		t.Error("expected SECRET_KEY in case-insensitive output")
	}
}

func TestGrepNoMatchPrintsMessage(t *testing.T) {
	r := setupGrepVault(t)
	var buf bytes.Buffer
	r.Stdout = &buf
	code := r.Run([]string{"grep", "zzznomatch"})
	if code != 0 {
		t.Fatalf("expected exit 0, got %d", code)
	}
	if !strings.Contains(buf.String(), "no matches") {
		t.Error("expected 'no matches' message")
	}
}

func TestGrepFailsWithoutInit(t *testing.T) {
	r := tempRunner(t)
	code := r.Run([]string{"grep", "anything"})
	if code == 0 {
		t.Error("expected non-zero exit without init")
	}
}

func TestGrepRequiresPattern(t *testing.T) {
	r := setupGrepVault(t)
	code := r.Run([]string{"grep"})
	if code == 0 {
		t.Error("expected non-zero exit with no pattern")
	}
}
