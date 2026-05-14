package cli

import (
	"testing"
)

func TestLoadSecretsFromContextReturnsSortedPairs(t *testing.T) {
	r := tempRunner(t)
	if err := r.Init(); err != nil {
		t.Fatalf("init: %v", err)
	}
	if err := r.Lock(); err != nil {
		t.Fatalf("lock: %v", err)
	}

	data := map[string]string{
		"ZEBRA": "z",
		"ALPHA": "a",
		"MANGO": "m",
	}
	if err := r.Vault.Write(r.Key, data); err != nil {
		t.Fatalf("write vault: %v", err)
	}

	pairs, err := loadSecretsFromRunner(r)
	if err != nil {
		t.Fatalf("loadSecrets: %v", err)
	}

	expected := []string{"ALPHA", "MANGO", "ZEBRA"}
	for i, kv := range pairs {
		if kv[0] != expected[i] {
			t.Errorf("index %d: want key %q, got %q", i, expected[i], kv[0])
		}
	}
}

func TestVaultPathFromRunner(t *testing.T) {
	r := tempRunner(t)
	path := r.Vault.Path()
	if path == "" {
		t.Error("expected non-empty vault path")
	}
}
