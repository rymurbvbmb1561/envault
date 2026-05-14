package cli

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/urfave/cli/v2"
)

func runWatch(c *cli.Context) error {
	if err := requireInit(c); err != nil {
		return err
	}
	if err := requireVault(c); err != nil {
		return err
	}

	output := c.String("output")
	exportPrefix := c.Bool("export")

	vaultPath := vaultPathFromContext(c)

	var lastMod time.Time

	writeOutput := func() error {
		secrets, err := loadSecretsFromContext(c)
		if err != nil {
			return fmt.Errorf("failed to decrypt vault: %w", err)
		}
		f, err := os.Create(output)
		if err != nil {
			return fmt.Errorf("failed to open output file: %w", err)
		}
		defer f.Close()
		for _, kv := range secrets {
			if exportPrefix {
				fmt.Fprintf(f, "export %s=%s\n", kv[0], kv[1])
			} else {
				fmt.Fprintf(f, "%s=%s\n", kv[0], kv[1])
			}
		}
		return nil
	}

	if err := writeOutput(); err != nil {
		return err
	}
	fmt.Fprintf(c.App.Writer, "watching %s → %s (Ctrl+C to stop)\n", vaultPath, output)

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-sig:
			fmt.Fprintln(c.App.Writer, "stopping watch")
			return nil
		case <-ticker.C:
			info, err := os.Stat(vaultPath)
			if err != nil {
				continue
			}
			if info.ModTime().After(lastMod) {
				lastMod = info.ModTime()
				if err := writeOutput(); err != nil {
					fmt.Fprintf(c.App.ErrWriter, "watch error: %v\n", err)
					continue
				}
				fmt.Fprintf(c.App.Writer, "[%s] vault changed, updated %s\n",
					time.Now().Format("15:04:05"), output)
			}
		}
	}
}
