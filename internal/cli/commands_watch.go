package cli

import "github.com/urfave/cli/v2"

func init() {
	registerCommand(watchCommand())
}

func watchCommand() *cli.Command {
	return &cli.Command{
		Name:  "watch",
		Usage: "watch the vault for changes and re-export to a .env file",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "output",
				Aliases: []string{"o"},
				Value:   ".env",
				Usage:   "output file to write decrypted secrets to",
			},
			&cli.BoolFlag{
				Name:  "export",
				Usage: "prefix each line with export",
			},
		},
		Action: runWatch,
	}
}
