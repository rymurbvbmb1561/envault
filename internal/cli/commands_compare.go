package cli

import (
	"github.com/urfave/cli/v2"
)

func init() {
	registerCommand(&cli.Command{
		Name:      "compare",
		Usage:     "compare secrets between two vault snapshots or a snapshot and the current vault",
		ArgsUsage: "<snapshot-a> [snapshot-b]",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "values",
				Aliases: []string{"v"},
				Usage:   "show values in comparison output",
			},
			&cli.BoolFlag{
				Name:  "keys-only",
				Usage: "only compare key names, ignore value changes",
			},
		},
		Action: func(c *cli.Context) error {
			return runCompare(c)
		},
	})
}
