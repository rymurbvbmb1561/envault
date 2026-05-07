package cli

import "os"

// readFileBytes is a thin wrapper around os.ReadFile so it can be swapped
// in tests if needed.
func readFileBytes(path string) ([]byte, error) {
	return os.ReadFile(path)
}
