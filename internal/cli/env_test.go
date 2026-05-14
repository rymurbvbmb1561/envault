package cli

import (
	"strings"
	"testing"
)

func TestEnvPrintsKeyValuePairs(t *testing.T) {
	r := tempRunner(t)
	if rc := runInit(r, nil); rc != 0 {
		t.Fatal("init failed")
	}
	if rc := runLock(r, []string{"API_KEY=secret", "DB_URL=postgres://localhost"}); rc != 0 {
		t.Fatal("lock failed")
	}

	out := captureStdout(t, func() {
		if rc := runEnv(r, nil); rc != 0 {
			t.Error("env returned non-zero")
		}
	})

	if !strings.Contains(out, "API_KEY='secret'") {
		t.Errorf("expected API_KEY in output, got: %s", out)
	}
	if !strings.Contains(out, "DB_URL='postgres://localhost'") {
		t.Errorf("expected DB_URL in output, got: %s", out)
	}
}

func TestEnvExportFlag(t *testing.T) {
	r := tempRunner(t)
	if rc := runInit(r, nil); rc != 0 {
		t.Fatal("init failed")
	}
	if rc := runLock(r, []string{"TOKEN=abc123"}); rc != 0 {
		t.Fatal("lock failed")
	}

	out := captureStdout(t, func() {
		if rc := runEnv(r, []string{"--export"}); rc != 0 {
			t.Error("env --export returned non-zero")
		}
	})

	if !strings.Contains(out, "export TOKEN='abc123'") {
		t.Errorf("expected 'export TOKEN=...' in output, got: %s", out)
	}
}

func TestEnvPrefixFilter(t *testing.T) {
	r := tempRunner(t)
	if rc := runInit(r, nil); rc != 0 {
		t.Fatal("init failed")
	}
	if rc := runLock(r, []string{"AWS_KEY=key", "AWS_SECRET=sec", "DB_PASS=pass"}); rc != 0 {
		t.Fatal("lock failed")
	}

	out := captureStdout(t, func() {
		if rc := runEnv(r, []string{"AWS"}); rc != 0 {
			t.Error("env AWS returned non-zero")
		}
	})

	if !strings.Contains(out, "AWS_KEY") {
		t.Errorf("expected AWS_KEY in output")
	}
	if !strings.Contains(out, "AWS_SECRET") {
		t.Errorf("expected AWS_SECRET in output")
	}
	if strings.Contains(out, "DB_PASS") {
		t.Errorf("DB_PASS should be filtered out")
	}
}

func TestEnvFailsWithoutInit(t *testing.T) {
	r := tempRunner(t)
	if rc := runEnv(r, nil); rc != 1 {
		t.Errorf("expected exit 1 without init, got %d", rc)
	}
}

func TestEnvFailsWithoutVault(t *testing.T) {
	r := tempRunner(t)
	if rc := runInit(r, nil); rc != 0 {
		t.Fatal("init failed")
	}
	if rc := runEnv(r, nil); rc != 1 {
		t.Errorf("expected exit 1 without vault, got %d", rc)
	}
}
