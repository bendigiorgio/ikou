package cmd

import (
	"os/exec"

	"github.com/urfave/cli/v2"
)

func GetDevCommand() *cli.Command {
	return &cli.Command{
		Name:  "dev",
		Usage: "Start the development server",
		Action: func(c *cli.Context) error {
			cmd := exec.Command("air")
			cmd.Dir = "../bin"
			err := cmd.Run()
			if err != nil {
				return err
			}
			return nil
		},
	}
}
