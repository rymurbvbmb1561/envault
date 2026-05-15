package cli

import (
	"fmt"
	"sort"
	"strings"

	"github.com/urfave/cli/v2"
)

func init() {
	registerCommand(pinCommand())
}

func pinCommand() *cli.Command {
	return &cli.Command{
		Name:  "pin",
		Usage: "pin or unpin keys to mark them as critical",
		Subcommands: []*cli.Command{
			{
				Name:      "add",
				Usage:     "pin one or more keys",
				ArgsUsage: "<KEY> [KEY...]",
				Action:    runPinAdd,
			},
			{
				Name:      "remove",
				Usage:     "unpin one or more keys",
				ArgsUsage: "<KEY> [KEY...]",
				Action:    runPinRemove,
			},
			{
				Name:   "list",
				Usage:  "list all pinned keys",
				Action: runPinList,
			},
		},
	}
}

func runPinAdd(c *cli.Context) error {
	if c.NArg() == 0 {
		return fmt.Errorf("at least one key name is required")
	}
	r, err := runnerFromContext(c)
	if err != nil {
		return err
	}
	if err := requiresInit(r); err != nil {
		return err
	}
	secrets, err := r.Vault.Read(r.Identity)
	if err != nil {
		return fmt.Errorf("could not read vault: %w", err)
	}
	for _, key := range c.Args().Slice() {
		key = strings.TrimSpace(key)
		if _, ok := secrets[key]; !ok {
			return fmt.Errorf("key %q not found in vault", key)
		}
		secrets["__pin__"+key] = "1"
		fmt.Fprintf(c.App.Writer, "pinned %s\n", key)
	}
	return r.Vault.Write(r.Recipient, secrets)
}

func runPinRemove(c *cli.Context) error {
	if c.NArg() == 0 {
		return fmt.Errorf("at least one key name is required")
	}
	r, err := runnerFromContext(c)
	if err != nil {
		return err
	}
	if err := requiresInit(r); err != nil {
		return err
	}
	secrets, err := r.Vault.Read(r.Identity)
	if err != nil {
		return fmt.Errorf("could not read vault: %w", err)
	}
	for _, key := range c.Args().Slice() {
		key = strings.TrimSpace(key)
		pinKey := "__pin__" + key
		if _, ok := secrets[pinKey]; !ok {
			return fmt.Errorf("key %q is not pinned", key)
		}
		delete(secrets, pinKey)
		fmt.Fprintf(c.App.Writer, "unpinned %s\n", key)
	}
	return r.Vault.Write(r.Recipient, secrets)
}

func runPinList(c *cli.Context) error {
	r, err := runnerFromContext(c)
	if err != nil {
		return err
	}
	if err := requiresInit(r); err != nil {
		return err
	}
	secrets, err := r.Vault.Read(r.Identity)
	if err != nil {
		return fmt.Errorf("could not read vault: %w", err)
	}
	var pinned []string
	for k := range secrets {
		if strings.HasPrefix(k, "__pin__") {
			pinned = append(pinned, strings.TrimPrefix(k, "__pin__"))
		}
	}
	if len(pinned) == 0 {
		fmt.Fprintln(c.App.Writer, "no pinned keys")
		return nil
	}
	sort.Strings(pinned)
	for _, k := range pinned {
		fmt.Fprintf(c.App.Writer, "📌 %s\n", k)
	}
	return nil
}
