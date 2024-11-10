package cmd

import (
	"github.com/urfave/cli/v2"
)

func GetServeCommand() *cli.Command {
	return &cli.Command{
		Name:  "serve",
		Usage: "Build Static Files",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "config",
				Aliases: []string{"c"},
				Usage:   "Path to the config file",
				Value:   "ikou.config.json",
			},
			&cli.BoolFlag{
				Name:    "api",
				Aliases: []string{"a"},
				Usage:   "Enable API routes",
				Value:   false,
			},
		},
		Action: func(c *cli.Context) error {
			//TODO: Implement serve command
			return nil
		},
	}
}
