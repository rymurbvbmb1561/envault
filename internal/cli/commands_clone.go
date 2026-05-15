package cli

import (
	"github.com/urfave/cli/v2"
)

func init() {
	additionalCommands = append(additionalCommands, &cli.Command{
		Name:      "clone",
		Usage:     "copy a secret to a new key",
		ArgsUsage: "<source> <destination>",
		Description: `Copies the value of an existing secret to a new key.

The source key must exist. The destination key must not already exist
unless --force is provided.

Example:
  envault clone DB_PASSWORD DB_PASSWORD_BACKUP`,
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "force",
				Aliases: []string{"f"},
				Usage:   "overwrite destination key if it already exists",
			},
		},
		Action: runClone,
	})
}
