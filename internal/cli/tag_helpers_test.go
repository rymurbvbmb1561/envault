package cli

import (
	"testing"
)

func TestRemoveTagsAllRemoved(t *testing.T) {
	result := removeTags([]string{"x", "y"}, []string{"x", "y"})
	if len(result) != 0 {
		t.Errorf("expected empty slice after removing all tags, got %v", result)
	}
}

func TestRemoveTagsPartial(t *testing.T) {
	result := removeTags([]string{"a", "b", "c"}, []string{"b"})
	if len(result) != 2 {
		t.Errorf("expected 2 tags, got %d: %v", len(result), result)
	}
	for _, t2 := range result {
		if t2 == "b" {
			t.Error("tag 'b' should have been removed")
		}
	}
}

func TestRemoveTagsNoop(t *testing.T) {
	result := removeTags([]string{"a", "b"}, []string{"z"})
	if len(result) != 2 {
		t.Errorf("expected 2 tags unchanged, got %d", len(result))
	}
}

func TestMergeTagsDeduplicates(t *testing.T) {
	result := mergeTags([]string{"foo", "bar"}, []string{"foo", "baz"})
	seen := map[string]int{}
	for _, tag := range result {
		seen[tag]++
	}
	for tag, count := range seen {
		if count > 1 {
			t.Errorf("tag %q appears %d times, expected 1", tag, count)
		}
	}
	if len(result) != 3 {
		t.Errorf("expected 3 unique tags, got %d: %v", len(result), result)
	}
}

func TestMergeTagsSorted(t *testing.T) {
	result := mergeTags([]string{"zebra"}, []string{"apple"})
	if result[0] != "apple" || result[1] != "zebra" {
		t.Errorf("expected sorted result, got %v", result)
	}
}

func TestParseTagsTrimsSpace(t *testing.T) {
	result := parseTags(" a , b , c ")
	if len(result) != 3 {
		t.Errorf("expected 3 tags after trimming, got %d: %v", len(result), result)
	}
}
