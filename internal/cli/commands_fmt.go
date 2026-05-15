package cli

import (
	"github.com/urfave/cli/v2"
)

func init() {
	registerCommand(&cli.Command{
		Name:  "fmt",
		Usage: "Sort and normalise keys in the vault",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "check",
				Usage: "Exit with a non-zero status if the vault is not already formatted",
			},
		},
		Action: func(c *cli.Context) error {
			return runFmt(c)
		},
	})
}
