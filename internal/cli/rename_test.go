package cli

import (
	"strings"
	"testing"
)

func TestRenameMovesKey(t *testing.T) {
	r := tempRunner(t)

	if rc := runInit(r, nil); rc != 0 {
		t.Fatal("init failed")
	}
	if rc := runAdd(r, []string{"OLD_KEY=hello"}); rc != 0 {
		t.Fatal("add failed")
	}

	r.Stdout.(*strings.Builder).Reset()
	if rc := runRename(r, []string{"OLD_KEY", "NEW_KEY"}); rc != 0 {
		t.Fatalf("rename returned non-zero: %s", r.Stderr.(*strings.Builder).String())
	}

	env, err := r.Vault.Read(r.Identity)
	if err != nil {
		t.Fatalf("read vault: %v", err)
	}
	if _, ok := env["OLD_KEY"]; ok {
		t.Error("OLD_KEY should have been removed")
	}
	if v, ok := env["NEW_KEY"]; !ok || v != "hello" {
		t.Errorf("expected NEW_KEY=hello, got %q (exists=%v)", v, ok)
	}
}

func TestRenameFailsWhenKeyMissing(t *testing.T) {
	r := tempRunner(t)
	runInit(r, nil)

	if rc := runRename(r, []string{"GHOST", "NEW_KEY"}); rc == 0 {
		t.Error("expected non-zero exit for missing key")
	}
	if !strings.Contains(r.Stderr.(*strings.Builder).String(), "not found") {
		t.Error("expected 'not found' in stderr")
	}
}

func TestRenameFailsWhenTargetExists(t *testing.T) {
	r := tempRunner(t)
	runInit(r, nil)
	runAdd(r, []string{"KEY_A=foo"})
	runAdd(r, []string{"KEY_B=bar"})

	if rc := runRename(r, []string{"KEY_A", "KEY_B"}); rc == 0 {
		t.Error("expected non-zero exit when target key already exists")
	}
	if !strings.Contains(r.Stderr.(*strings.Builder).String(), "already exists") {
		t.Error("expected 'already exists' in stderr")
	}
}

func TestRenameFailsWithoutInit(t *testing.T) {
	r := tempRunner(t)

	if rc := runRename(r, []string{"A", "B"}); rc == 0 {
		t.Error("expected non-zero exit without init")
	}
}

func TestRenameRequiresTwoArgs(t *testing.T) {
	r := tempRunner(t)
	runInit(r, nil)

	for _, args := range [][]string{{}, {"ONLY_ONE"}} {
		if rc := runRename(r, args); rc == 0 {
			t.Errorf("expected non-zero for args %v", args)
		}
	}
}

func TestRenameSameKeyReturnsError(t *testing.T) {
	r := tempRunner(t)
	runInit(r, nil)
	runAdd(r, []string{"MY_KEY=value"})

	if rc := runRename(r, []string{"MY_KEY", "MY_KEY"}); rc == 0 {
		t.Error("expected non-zero when old and new names are identical")
	}
}
