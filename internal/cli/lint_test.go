package cli

import (
	"strings"
	"testing"
)

func TestLintNoIssues(t *testing.T) {
	r := tempRunner(t)
	if code := runInit(r, nil); code != 0 {
		t.Fatalf("init failed: %s", r.Stderr)
	}
	if code := runLock(r, nil); code != 0 {
		t.Fatalf("lock failed: %s", r.Stderr)
	}
	if code := runAdd(r, []string{"DATABASE_URL=postgres://localhost/db"}); code != 0 {
		t.Fatalf("add failed: %s", r.Stderr)
	}
	if code := runAdd(r, []string{"API_KEY=secret"}); code != 0 {
		t.Fatalf("add failed: %s", r.Stderr)
	}

	code := runLint(r, nil)
	if code != 0 {
		t.Fatalf("expected exit 0, got %d; stderr: %s", code, r.Stderr)
	}
	if !strings.Contains(r.Stdout.String(), "no lint issues found") {
		t.Errorf("expected no-issue message, got: %s", r.Stdout)
	}
}

func TestLintDetectsLowercaseKey(t *testing.T) {
	r := tempRunner(t)
	if code := runInit(r, nil); code != 0 {
		t.Fatalf("init failed")
	}
	if code := runLock(r, nil); code != 0 {
		t.Fatalf("lock failed")
	}
	if code := runAdd(r, []string{"database_url=postgres://localhost/db"}); code != 0 {
		t.Fatalf("add failed")
	}

	code := runLint(r, nil)
	if code != 1 {
		t.Fatalf("expected exit 1 for lowercase key, got %d", code)
	}
	if !strings.Contains(r.Stdout.String(), "lowercase") {
		t.Errorf("expected lowercase warning, got: %s", r.Stdout)
	}
}

func TestLintDetectsTrailingUnderscore(t *testing.T) {
	r := tempRunner(t)
	if code := runInit(r, nil); code != 0 {
		t.Fatalf("init failed")
	}
	if code := runLock(r, nil); code != 0 {
		t.Fatalf("lock failed")
	}
	if code := runAdd(r, []string{"API_KEY_=secret"}); code != 0 {
		t.Fatalf("add failed")
	}

	code := runLint(r, nil)
	if code != 1 {
		t.Fatalf("expected exit 1 for trailing underscore, got %d", code)
	}
	if !strings.Contains(r.Stdout.String(), "ends with underscore") {
		t.Errorf("expected trailing underscore warning, got: %s", r.Stdout)
	}
}

func TestLintFailsWithoutInit(t *testing.T) {
	r := tempRunner(t)
	code := runLint(r, nil)
	if code != 1 {
		t.Fatalf("expected exit 1 without init, got %d", code)
	}
}

func TestLintFailsWithoutVault(t *testing.T) {
	r := tempRunner(t)
	if code := runInit(r, nil); code != 0 {
		t.Fatalf("init failed")
	}
	code := runLint(r, nil)
	if code != 1 {
		t.Fatalf("expected exit 1 without vault, got %d", code)
	}
}
