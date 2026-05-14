package cli

import (
	"strings"
	"testing"
)

func TestAuditNoIssues(t *testing.T) {
	r := tempRunner(t)
	if rc := r.Run([]string{"init"}); rc != 0 {
		t.Fatalf("init failed")
	}
	if rc := r.Run([]string{"add", "API_KEY=supersecurevalue123"}); rc != 0 {
		t.Fatalf("add failed")
	}
	if rc := r.Run([]string{"lock"}); rc != 0 {
		t.Fatalf("lock failed")
	}
	if rc := r.Run([]string{"audit"}); rc != 0 {
		t.Errorf("expected exit 0, got non-zero; stderr: %s", r.Stderr)
	}
	if !strings.Contains(r.Stdout.String(), "no issues found") {
		t.Errorf("expected 'no issues found', got: %s", r.Stdout)
	}
}

func TestAuditDetectsEmptyValue(t *testing.T) {
	r := tempRunner(t)
	if rc := r.Run([]string{"init"}); rc != 0 {
		t.Fatalf("init failed")
	}
	if rc := r.Run([]string{"add", "EMPTY_KEY="}); rc != 0 {
		t.Fatalf("add failed")
	}
	if rc := r.Run([]string{"lock"}); rc != 0 {
		t.Fatalf("lock failed")
	}
	rc := r.Run([]string{"audit"})
	if rc == 0 {
		t.Errorf("expected non-zero exit for audit with issues")
	}
	out := r.Stdout.String()
	if !strings.Contains(out, "EMPTY_KEY") {
		t.Errorf("expected EMPTY_KEY in output, got: %s", out)
	}
	if !strings.Contains(out, "empty value") {
		t.Errorf("expected 'empty value' in output, got: %s", out)
	}
}

func TestAuditDetectsPlaceholder(t *testing.T) {
	r := tempRunner(t)
	if rc := r.Run([]string{"init"}); rc != 0 {
		t.Fatalf("init failed")
	}
	if rc := r.Run([]string{"add", "DB_PASS=changeme"}); rc != 0 {
		t.Fatalf("add failed")
	}
	if rc := r.Run([]string{"lock"}); rc != 0 {
		t.Fatalf("lock failed")
	}
	rc := r.Run([]string{"audit"})
	if rc == 0 {
		t.Errorf("expected non-zero exit")
	}
	out := r.Stdout.String()
	if !strings.Contains(out, "placeholder value") {
		t.Errorf("expected placeholder warning, got: %s", out)
	}
}

func TestAuditDetectsDuplicateValues(t *testing.T) {
	r := tempRunner(t)
	if rc := r.Run([]string{"init"}); rc != 0 {
		t.Fatalf("init failed")
	}
	if rc := r.Run([]string{"add", "KEY_A=sharedvalue"}); rc != 0 {
		t.Fatalf("add KEY_A failed")
	}
	if rc := r.Run([]string{"add", "KEY_B=sharedvalue"}); rc != 0 {
		t.Fatalf("add KEY_B failed")
	}
	if rc := r.Run([]string{"lock"}); rc != 0 {
		t.Fatalf("lock failed")
	}
	rc := r.Run([]string{"audit"})
	if rc == 0 {
		t.Errorf("expected non-zero exit for duplicate values")
	}
	out := r.Stdout.String()
	if !strings.Contains(out, "duplicate value") {
		t.Errorf("expected 'duplicate value' in output, got: %s", out)
	}
}

func TestAuditFailsWithoutInit(t *testing.T) {
	r := tempRunner(t)
	rc := r.Run([]string{"audit"})
	if rc == 0 {
		t.Errorf("expected non-zero exit without init")
	}
}

func TestAuditFailsWithoutVault(t *testing.T) {
	r := tempRunner(t)
	if rc := r.Run([]string{"init"}); rc != 0 {
		t.Fatalf("init failed")
	}
	rc := r.Run([]string{"audit"})
	if rc == 0 {
		t.Errorf("expected non-zero exit without vault")
	}
}
