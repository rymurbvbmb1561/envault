package cli

import (
	"bytes"
	"strings"
	"testing"
)

func setupSearchVault(t *testing.T) *Runner {
	t.Helper()
	r := tempRunner(t)
	if code := runInit(r, nil); code != 0 {
		t.Fatal("init failed")
	}
	if code := runLock(r, nil); code != 0 {
		t.Fatal("lock failed")
	}
	for _, kv := range []string{"DB_HOST=localhost", "DB_PORT=5432", "API_KEY=secret123", "APP_ENV=production"} {
		if code := runAdd(r, []string{kv}); code != 0 {
			t.Fatalf("add %s failed", kv)
		}
	}
	return r
}

func TestSearchMatchesKeyPrefix(t *testing.T) {
	r := setupSearchVault(t)
	r.Stdout.(*bytes.Buffer).Reset()
	code := runSearch(r, []string{"DB_"})
	if code != 0 {
		t.Fatalf("search returned %d", code)
	}
	out := r.Stdout.(*bytes.Buffer).String()
	if !strings.Contains(out, "DB_HOST") || !strings.Contains(out, "DB_PORT") {
		t.Errorf("expected DB keys in output, got: %s", out)
	}
	if strings.Contains(out, "API_KEY") {
		t.Errorf("unexpected API_KEY in output: %s", out)
	}
}

func TestSearchNoMatchPrintsMessage(t *testing.T) {
	r := setupSearchVault(t)
	r.Stdout.(*bytes.Buffer).Reset()
	code := runSearch(r, []string{"ZZZNOMATCH"})
	if code != 0 {
		t.Fatalf("expected zero exit, got %d", code)
	}
	out := r.Stdout.(*bytes.Buffer).String()
	if !strings.Contains(out, "no keys matched") {
		t.Errorf("expected no-match message, got: %s", out)
	}
}

func TestSearchWithValuesFlag(t *testing.T) {
	r := setupSearchVault(t)
	r.Stdout.(*bytes.Buffer).Reset()
	code := runSearch(r, []string{"secret", "--values"})
	if code != 0 {
		t.Fatalf("search returned %d", code)
	}
	out := r.Stdout.(*bytes.Buffer).String()
	if !strings.Contains(out, "API_KEY") {
		t.Errorf("expected API_KEY matched by value, got: %s", out)
	}
}

func TestSearchFailsWithoutInit(t *testing.T) {
	r := tempRunner(t)
	code := runSearch(r, []string{"DB"})
	if code == 0 {
		t.Fatal("expected non-zero exit without init")
	}
}

func TestSearchNoArgsReturnsError(t *testing.T) {
	r := tempRunner(t)
	code := runSearch(r, []string{})
	if code == 0 {
		t.Fatal("expected non-zero exit with no args")
	}
	if !strings.Contains(r.Stderr.(*bytes.Buffer).String(), "usage") {
		t.Errorf("expected usage hint")
	}
}
