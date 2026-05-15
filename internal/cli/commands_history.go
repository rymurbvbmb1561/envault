package cli

import (
	"github.com/urfave/cli/v2"
)

func init() {
	additionalCommands = append(additionalCommands, historyCommand)
}

var historyCommand = &cli.Command{
	Name:  "history",
	Usage: "Show a changelog of secret key mutations",
	Flags: []cli.Flag{
		&cli.IntFlag{
			Name:    "limit",
			Aliases: []string{"n"},
			Value:   20,
			Usage:   "Maximum number of entries to show",
		},
		&cli.StringFlag{
			Name:    "key",
			Aliases: []string{"k"},
			Usage:   "Filter history to a specific key",
		},
	},
	Action: runHistory,
}
