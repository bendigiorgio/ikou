package main

import (
	"os"
	"time"

	"github.com/bendigiorgio/ikou/internal/cmd"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:     "ikou",
		Version:  "0.0.1",
		Compiled: time.Now(),
		Authors: []*cli.Author{
			{
				Name:  "Ben Di Giorgio",
				Email: "bendigiorgio@gmail.com",
			},
		},
		HelpName: "ikou",
		Usage:    "A React Meta Framework",
		Commands: []*cli.Command{
			cmd.GetRunCommand(),
			cmd.GetDevCommand(),
			cmd.GetBuildCommand(),
		},
	}

	if err := app.Run(os.Args); err != nil {
		panic(err)
	}
}
