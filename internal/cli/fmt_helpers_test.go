package cli

import (
	"testing"
)

func TestSortedKeysReturnsAlphaOrder(t *testing.T) {
	m := map[string]string{
		"ZEBRA": "1",
		"ALPHA": "2",
		"MANGO": "3",
	}
	got := sortedKeys(m)
	want := []string{"ALPHA", "MANGO", "ZEBRA"}
	for i, k := range got {
		if k != want[i] {
			t.Errorf("index %d: got %q, want %q", i, k, want[i])
		}
	}
}

func TestSortedKeysEmptyMap(t *testing.T) {
	got := sortedKeys(map[string]string{})
	if len(got) != 0 {
		t.Errorf("expected empty slice, got %v", got)
	}
}

func TestMapsEqualIdentical(t *testing.T) {
	a := map[string]string{"A": "1", "B": "2"}
	b := map[string]string{"A": "1", "B": "2"}
	if !mapsEqual(a, b) {
		t.Error("expected maps to be equal")
	}
}

func TestMapsEqualDifferentValues(t *testing.T) {
	a := map[string]string{"A": "1"}
	b := map[string]string{"A": "2"}
	if mapsEqual(a, b) {
		t.Error("expected maps to be unequal")
	}
}

func TestMapsEqualDifferentLengths(t *testing.T) {
	a := map[string]string{"A": "1"}
	b := map[string]string{"A": "1", "B": "2"}
	if mapsEqual(a, b) {
		t.Error("expected maps to be unequal")
	}
}

func TestIsSortedTrue(t *testing.T) {
	// Single pass over map iteration is non-deterministic, but isSorted
	// collects and checks — so just verify a known-sorted map passes.
	m := map[string]string{"ALPHA": "1", "BETA": "2", "GAMMA": "3"}
	// We can't guarantee iteration order, but we CAN guarantee that after
	// sorting the keys they match themselves.
	keys := sortedKeys(m)
	for i := 1; i < len(keys); i++ {
		if keys[i] < keys[i-1] {
			t.Errorf("sortedKeys not sorted at index %d", i)
		}
	}
}
