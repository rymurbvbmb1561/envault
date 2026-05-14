package cli

import (
	"strings"
	"testing"
)

func TestValidateNoIssues(t *testing.T) {
	r := tempRunner(t)

	r.Commands["init"](r, nil)
	r.Commands["lock"](r, nil)

	r.Commands["add"](r, []string{"API_KEY=secret123"})
	r.Commands["add"](r, []string{"DATABASE_URL=postgres://localhost/db"})

	code := r.Commands["validate"](r, nil)
	if code != 0 {
		t.Fatalf("expected exit 0, got %d: %s", code, r.Stdout)
	}
	if !strings.Contains(r.Stdout.String(), "✓") {
		t.Errorf("expected success marker in output, got: %s", r.Stdout)
	}
}

func TestValidateDetectsLowercaseKey(t *testing.T) {
	issues := validateSecrets(map[string]string{
		"bad_key": "value",
	})
	if len(issues) == 0 {
		t.Fatal("expected an issue for lowercase key")
	}
	if issues[0].Key != "bad_key" {
		t.Errorf("unexpected key: %s", issues[0].Key)
	}
}

func TestValidateDetectsEmptyValue(t *testing.T) {
	issues := validateSecrets(map[string]string{
		"VALID_KEY": "   ",
	})
	if len(issues) == 0 {
		t.Fatal("expected an issue for empty value")
	}
	if !strings.Contains(issues[0].Message, "empty") {
		t.Errorf("expected empty value message, got: %s", issues[0].Message)
	}
}

func TestValidateDetectsNewlineInValue(t *testing.T) {
	issues := validateSecrets(map[string]string{
		"MY_KEY": "line1\nline2",
	})
	if len(issues) == 0 {
		t.Fatal("expected an issue for newline in value")
	}
	if !strings.Contains(issues[0].Message, "newline") {
		t.Errorf("expected newline message, got: %s", issues[0].Message)
	}
}

func TestValidateFailsWithoutInit(t *testing.T) {
	r := tempRunner(t)
	code := r.Commands["validate"](r, nil)
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}
}

func TestValidateFailsWithoutVault(t *testing.T) {
	r := tempRunner(t)
	r.Commands["init"](r, nil)
	code := r.Commands["validate"](r, nil)
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}
}
