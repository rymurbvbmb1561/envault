package cli

import (
	"fmt"
	"os"
)

// readFileBytes is a thin wrapper around os.ReadFile so it can be swapped
// in tests if needed.
func readFileBytes(path string) ([]byte, error) {
	return os.ReadFile(path)
}

// diffLines returns a simple line-level diff between two byte slices,
// representing added lines with "+" and removed lines with "-".
// It is intended for human-readable output only, not patch generation.
func diffLines(oldBytes, newBytes []byte) string {
	oldLines := splitLines(oldBytes)
	newLines := splitLines(newBytes)

	var result string
	for _, line := range oldLines {
		result += fmt.Sprintf("- %s\n", line)
	}
	for _, line := range newLines {
		result += fmt.Sprintf("+ %s\n", line)
	}
	return result
}

// splitLines splits a byte slice into a slice of strings by newline.
func splitLines(data []byte) []string {
	if len(data) == 0 {
		return nil
	}
	var lines []string
	start := 0
	for i, b := range data {
		if b == '\n' {
			lines = append(lines, string(data[start:i]))
			start = i + 1
		}
	}
	if start < len(data) {
		lines = append(lines, string(data[start:]))
	}
	return lines
}
