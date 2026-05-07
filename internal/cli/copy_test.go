package cli

import (
	"bytes"
	"runtime"
	"strings"
	"testing"
)

func TestCopyFailsWithoutInit(t *testing.T) {
	r := tempRunner(t)
	code := runCopy(r, []string{"SECRET_KEY"})
	if code == 0 {
		t.Fatal("expected non-zero exit without init")
	}
	if !strings.Contains(r.Stderr.(*bytes.Buffer).String(), "not initialised") {
		t.Errorf("expected 'not initialised' error, got: %s", r.Stderr.(*bytes.Buffer).String())
	}
}

func TestCopyFailsWithoutVault(t *testing.T) {
	r := tempRunner(t)
	if code := runInit(r, nil); code != 0 {
		t.Fatal("init failed")
	}
	code := runCopy(r, []string{"MISSING_KEY"})
	if code == 0 {
		t.Fatal("expected non-zero exit without vault")
	}
}

func TestCopyMissingKeyReturnsError(t *testing.T) {
	r := tempRunner(t)
	if code := runInit(r, nil); code != 0 {
		t.Fatal("init failed")
	}
	if code := runLock(r, nil); code != 0 {
		t.Fatal("lock failed")
	}
	if code := runAdd(r, []string{"FOO=bar"}); code != 0 {
		t.Fatal("add failed")
	}
	code := runCopy(r, []string{"DOES_NOT_EXIST"})
	if code == 0 {
		t.Fatal("expected non-zero exit for missing key")
	}
	if !strings.Contains(r.Stderr.(*bytes.Buffer).String(), "not found") {
		t.Errorf("expected 'not found' in error, got: %s", r.Stderr.(*bytes.Buffer).String())
	}
}

func TestCopyNoArgsReturnsError(t *testing.T) {
	r := tempRunner(t)
	code := runCopy(r, []string{})
	if code == 0 {
		t.Fatal("expected non-zero exit with no args")
	}
	if !strings.Contains(r.Stderr.(*bytes.Buffer).String(), "usage") {
		t.Errorf("expected usage hint, got: %s", r.Stderr.(*bytes.Buffer).String())
	}
}

func TestCopyToClipboardUnsupportedPlatform(t *testing.T) {
	if runtime.GOOS == "darwin" || runtime.GOOS == "linux" || runtime.GOOS == "windows" {
		t.Skip("skipping unsupported platform test on supported OS")
	}
	err := copyToClipboard("test")
	if err == nil {
		t.Fatal("expected error on unsupported platform")
	}
}
