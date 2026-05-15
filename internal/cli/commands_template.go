package cli

import (
	"github.com/urfave/cli/v2"
)

func init() {
	extraCommands = append(extraCommands, templateCommand())
}

func templateCommand() *cli.Command {
	return &cli.Command{
		Name:      "template",
		Usage:     "render a template file substituting secret values",
		ArgsUsage: "<template-file>",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "output",
				Aliases: []string{"o"},
				Usage:   "write rendered output to `FILE` instead of stdout",
			},
			&cli.StringFlag{
				Name:  "prefix",
				Usage: "only expose secrets whose key starts with `PREFIX`",
			},
		},
		Action: func(c *cli.Context) error {
			r, err := runnerFromContext(c)
			if err != nil {
				return err
			}
			return runTemplate(c, r)
		},
	}
}
