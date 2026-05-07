package cli

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeDotenv(t *testing.T, dir, content string) string {
	t.Helper()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0600); err != nil {
		t.Fatalf("writeDotenv: %v", err)
	}
	return p
}

func TestDiffNoChanges(t *testing.T) {
	r, dir := tempRunner(t)
	if code := r.Run([]string{"init"}); code != 0 {
		t.Fatalf("init failed")
	}
	if code := r.Run([]string{"add", "FOO=bar"}); code != 0 {
		t.Fatalf("add failed")
	}
	if code := r.Run([]string{"lock"}); code != 0 {
		t.Fatalf("lock failed")
	}

	dotenv := writeDotenv(t, dir, "FOO=bar\n")
	out := &strings.Builder{}
	r.Stdout = out
	if code := r.Run([]string{"diff", dotenv}); code != 0 {
		t.Fatalf("diff returned non-zero: %s", r.Stderr)
	}
	if !strings.Contains(out.String(), "no differences") {
		t.Errorf("expected 'no differences', got %q", out.String())
	}
}

func TestDiffDetectsAddition(t *testing.T) {
	r, dir := tempRunner(t)
	if code := r.Run([]string{"init"}); code != 0 {
		t.Fatalf("init failed")
	}
	if code := r.Run([]string{"add", "FOO=bar"}); code != 0 {
		t.Fatalf("add failed")
	}
	if code := r.Run([]string{"lock"}); code != 0 {
		t.Fatalf("lock failed")
	}

	dotenv := writeDotenv(t, dir, "FOO=bar\nNEW_KEY=secret\n")
	out := &strings.Builder{}
	r.Stdout = out
	if code := r.Run([]string{"diff", dotenv}); code != 0 {
		t.Fatalf("diff returned non-zero")
	}
	if !strings.Contains(out.String(), "+ NEW_KEY") {
		t.Errorf("expected addition marker, got %q", out.String())
	}
}

func TestDiffDetectsRemoval(t *testing.T) {
	r, dir := tempRunner(t)
	if code := r.Run([]string{"init"}); code != 0 {
		t.Fatalf("init failed")
	}
	if code := r.Run([]string{"add", "FOO=bar"}); code != 0 {
		t.Fatalf("add failed")
	}
	if code := r.Run([]string{"add", "GONE=yes"}); code != 0 {
		t.Fatalf("add failed")
	}
	if code := r.Run([]string{"lock"}); code != 0 {
		t.Fatalf("lock failed")
	}

	dotenv := writeDotenv(t, dir, "FOO=bar\n")
	out := &strings.Builder{}
	r.Stdout = out
	if code := r.Run([]string{"diff", dotenv}); code != 0 {
		t.Fatalf("diff returned non-zero")
	}
	if !strings.Contains(out.String(), "- GONE") {
		t.Errorf("expected removal marker, got %q", out.String())
	}
}

func TestDiffDetectsChange(t *testing.T) {
	r, dir := tempRunner(t)
	if code := r.Run([]string{"init"}); code != 0 {
		t.Fatalf("init failed")
	}
	if code := r.Run([]string{"add", "FOO=original"}); code != 0 {
		t.Fatalf("add failed")
	}
	if code := r.Run([]string{"lock"}); code != 0 {
		t.Fatalf("lock failed")
	}

	dotenv := writeDotenv(t, dir, "FOO=changed\n")
	out := &strings.Builder{}
	r.Stdout = out
	if code := r.Run([]string{"diff", dotenv}); code != 0 {
		t.Fatalf("diff returned non-zero")
	}
	if !strings.Contains(out.String(), "~ FOO") {
		t.Errorf("expected change marker, got %q", out.String())
	}
}

func TestDiffFailsWithoutInit(t *testing.T) {
	r, _ := tempRunner(t)
	if code := r.Run([]string{"diff"}); code == 0 {
		t.Error("expected non-zero exit without init")
	}
}

func TestDiffFailsWithoutVault(t *testing.T) {
	r, _ := tempRunner(t)
	if code := r.Run([]string{"init"}); code != 0 {
		t.Fatalf("init failed")
	}
	if code := r.Run([]string{"diff"}); code == 0 {
		t.Error("expected non-zero exit without vault")
	}
}
