package cli

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeTemplate(t *testing.T, dir, name, content string) string {
	t.Helper()
	p := filepath.Join(dir, name)
	if err := os.WriteFile(p, []byte(content), 0644); err != nil {
		t.Fatalf("write template: %v", err)
	}
	return p
}

func TestTemplateRendersSecrets(t *testing.T) {
	r := tempRunner(t)
	run(t, r, "init")
	run(t, r, "add", "APP_HOST=localhost")
	run(t, r, "add", "APP_PORT=8080")

	tmplPath := writeTemplate(t, r.Config.ProjectRoot, "app.tmpl",
		"host={{index . \"APP_HOST\"}} port={{index . \"APP_PORT\"}}")

	out := captureOutput(t, r, "template", tmplPath)
	if !strings.Contains(out, "host=localhost") {
		t.Errorf("expected host=localhost in output, got: %s", out)
	}
	if !strings.Contains(out, "port=8080") {
		t.Errorf("expected port=8080 in output, got: %s", out)
	}
}

func TestTemplateWritesToOutputFile(t *testing.T) {
	r := tempRunner(t)
	run(t, r, "init")
	run(t, r, "add", "DB_URL=postgres://localhost/mydb")

	tmplPath := writeTemplate(t, r.Config.ProjectRoot, "db.tmpl",
		"DATABASE_URL={{index . \"DB_URL\"}}")
	outPath := filepath.Join(r.Config.ProjectRoot, "rendered.env")

	run(t, r, "template", "--output", outPath, tmplPath)

	data, err := os.ReadFile(outPath)
	if err != nil {
		t.Fatalf("read output file: %v", err)
	}
	if !strings.Contains(string(data), "postgres://localhost/mydb") {
		t.Errorf("expected DB_URL value in file, got: %s", data)
	}
}

func TestTemplateFailsWithoutInit(t *testing.T) {
	r := tempRunner(t)
	tmplPath := writeTemplate(t, t.TempDir(), "x.tmpl", "hello")
	if err := runErr(r, "template", tmplPath); err == nil {
		t.Fatal("expected error when not initialised")
	}
}

func TestTemplateFailsWithoutVault(t *testing.T) {
	r := tempRunner(t)
	run(t, r, "init")
	if err := os.Remove(r.Config.VaultPath); err != nil {
		t.Fatalf("remove vault: %v", err)
	}
	tmplPath := writeTemplate(t, r.Config.ProjectRoot, "x.tmpl", "hello")
	if err := runErr(r, "template", tmplPath); err == nil {
		t.Fatal("expected error when vault missing")
	}
}

func TestTemplateFailsWithMissingKey(t *testing.T) {
	r := tempRunner(t)
	run(t, r, "init")
	run(t, r, "add", "EXISTING=val")

	tmplPath := writeTemplate(t, r.Config.ProjectRoot, "bad.tmpl",
		"{{index . \"MISSING_KEY\"}}")
	if err := runErr(r, "template", tmplPath); err == nil {
		t.Fatal("expected error for missing key in template")
	}
}

func TestTemplateNoArgsReturnsError(t *testing.T) {
	r := tempRunner(t)
	run(t, r, "init")
	if err := runErr(r, "template"); err == nil {
		t.Fatal("expected error when no template file provided")
	}
}
