package cli_test

import (
	"os"
	"testing"

	"github.com/user/envault/internal/cli"
)

func TestRunUnknownCommandReturnsOne(t *testing.T) {
	old := os.Args
	defer func() { os.Args = old }()

	os.Args = []string{"envault", "doesnotexist"}

	if code := cli.Run("test"); code != 1 {
		t.Errorf("expected exit code 1, got %d", code)
	}
}

func TestRunNoArgsReturnsOne(t *testing.T) {
	old := os.Args
	defer func() { os.Args = old }()

	os.Args = []string{"envault"}

	if code := cli.Run("test"); code != 1 {
		t.Errorf("expected exit code 1, got %d", code)
	}
}

func TestRunVersionReturnsZero(t *testing.T) {
	dir := t.TempDir()
	old := os.Args
	oldDir, _ := os.Getwd()
	defer func() {
		os.Args = old
		os.Chdir(oldDir)
	}()

	os.Chdir(dir)
	os.Args = []string{"envault", "version"}

	if code := cli.Run("test"); code != 0 {
		t.Errorf("expected exit code 0, got %d", code)
	}
}

func TestRunExportFormatFlagParsed(t *testing.T) {
	dir := t.TempDir()
	old := os.Args
	oldDir, _ := os.Getwd()
	defer func() {
		os.Args = old
		os.Chdir(oldDir)
	}()

	os.Chdir(dir)
	os.Args = []string{"envault", "export", "--format", "dotenv"}

	// No vault / no init — we just want the flag parsing not to panic.
	code := cli.Run("test")
	if code == 0 {
		t.Error("expected non-zero exit because project is not initialised")
	}
}
