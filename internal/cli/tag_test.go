package cli

import (
	"strings"
	"testing"
)

func TestTagAddAndList(t *testing.T) {
	r := tempRunner(t)
	if code := runInit(r, nil); code != 0 {
		t.Fatal("init failed")
	}
	if code := runAdd(r, []string{"DATABASE_URL=postgres://localhost/db"}); code != 0 {
		t.Fatal("add failed")
	}

	code := runTag(r, []string{"add", "DATABASE_URL", "production", "ci"})
	if code != 0 {
		t.Fatalf("tag add returned %d", code)
	}

	r.Stdout.Reset()
	code = runTag(r, []string{"list", "DATABASE_URL"})
	if code != 0 {
		t.Fatalf("tag list returned %d", code)
	}
	out := r.Stdout.String()
	if !strings.Contains(out, "ci") || !strings.Contains(out, "production") {
		t.Errorf("expected tags in output, got: %s", out)
	}
}

func TestTagRemove(t *testing.T) {
	r := tempRunner(t)
	if code := runInit(r, nil); code != 0 {
		t.Fatal("init failed")
	}
	if code := runAdd(r, []string{"API_KEY=secret"}); code != 0 {
		t.Fatal("add failed")
	}
	runTag(r, []string{"add", "API_KEY", "production", "deprecated"})

	code := runTag(r, []string{"remove", "API_KEY", "deprecated"})
	if code != 0 {
		t.Fatalf("tag remove returned %d", code)
	}

	r.Stdout.Reset()
	runTag(r, []string{"list", "API_KEY"})
	out := r.Stdout.String()
	if strings.Contains(out, "deprecated") {
		t.Errorf("expected deprecated to be removed, got: %s", out)
	}
	if !strings.Contains(out, "production") {
		t.Errorf("expected production to remain, got: %s", out)
	}
}

func TestTagListNoTags(t *testing.T) {
	r := tempRunner(t)
	runInit(r, nil)
	runAdd(r, []string{"PLAIN_KEY=value"})

	r.Stdout.Reset()
	code := runTag(r, []string{"list", "PLAIN_KEY"})
	if code != 0 {
		t.Fatalf("expected 0, got %d", code)
	}
	if !strings.Contains(r.Stdout.String(), "no tags") {
		t.Errorf("expected 'no tags' message, got: %s", r.Stdout.String())
	}
}

func TestTagFailsWithoutInit(t *testing.T) {
	r := tempRunner(t)
	code := runTag(r, []string{"list", "SOME_KEY"})
	if code != 1 {
		t.Errorf("expected exit 1 without init, got %d", code)
	}
}

func TestTagMissingKeyReturnsError(t *testing.T) {
	r := tempRunner(t)
	runInit(r, nil)

	code := runTag(r, []string{"add", "NONEXISTENT", "mytag"})
	if code != 1 {
		t.Errorf("expected exit 1 for missing key, got %d", code)
	}
}

func TestTagUnknownSubcommand(t *testing.T) {
	r := tempRunner(t)
	runInit(r, nil)
	runAdd(r, []string{"FOO=bar"})

	code := runTag(r, []string{"frobnicate", "FOO"})
	if code != 1 {
		t.Errorf("expected exit 1 for unknown subcommand, got %d", code)
	}
}

func TestParseTags(t *testing.T) {
	got := parseTags("a,b,c")
	if len(got) != 3 {
		t.Errorf("expected 3 tags, got %d", len(got))
	}
	if len(parseTags("")) != 0 {
		t.Error("expected empty slice for empty string")
	}
}

func TestMergeTags(t *testing.T) {
	result := mergeTags([]string{"a", "b"}, []string{"b", "c"})
	if len(result) != 3 {
		t.Errorf("expected 3 unique tags, got %d: %v", len(result), result)
	}
}
