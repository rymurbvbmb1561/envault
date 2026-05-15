package cli

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"text/template"

	"github.com/urfave/cli/v2"
)

func runTemplate(c *cli.Context, r *Runner) error {
	if err := requireInit(r); err != nil {
		return err
	}
	if err := requireVault(r); err != nil {
		return err
	}
	if c.NArg() < 1 {
		return fmt.Errorf("usage: envault template <template-file>")
	}

	tmplPath := c.Args().First()
	tmplBytes, err := os.ReadFile(tmplPath)
	if err != nil {
		return fmt.Errorf("read template: %w", err)
	}

	secrets, err := r.Vault.Read(r.KeyStore)
	if err != nil {
		return fmt.Errorf("unlock vault: %w", err)
	}

	prefix := c.String("prefix")
	data := make(map[string]string, len(secrets))
	for k, v := range secrets {
		if prefix != "" && !strings.HasPrefix(k, prefix) {
			continue
		}
		data[k] = v
	}

	tmpl, err := template.New("envault").Option("missingkey=error").Parse(string(tmplBytes))
	if err != nil {
		return fmt.Errorf("parse template: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return fmt.Errorf("render template: %w", err)
	}

	outPath := c.String("output")
	if outPath == "" {
		_, err = fmt.Fprint(c.App.Writer, buf.String())
		return err
	}

	if err := os.WriteFile(outPath, buf.Bytes(), 0600); err != nil {
		return fmt.Errorf("write output file: %w", err)
	}
	fmt.Fprintf(c.App.Writer, "rendered template written to %s\n", outPath)
	return nil
}
