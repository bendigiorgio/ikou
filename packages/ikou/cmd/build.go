package cmd

import (
	"github.com/urfave/cli/v2"
)

func GetBuildCommand() *cli.Command {
	return &cli.Command{
		Name:  "run",
		Usage: "Build Static Files",
		Action: func(c *cli.Context) error {
			//TODO: Implement build command
			return nil
		},
	}
}
