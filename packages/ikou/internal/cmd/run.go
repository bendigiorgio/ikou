package cmd

import (
	"github.com/bendigiorgio/ikou/internal/app"
	"github.com/bendigiorgio/ikou/internal/app/react"
	"github.com/bendigiorgio/ikou/internal/app/utils"
	"github.com/urfave/cli/v2"
)

func GetRunCommand() *cli.Command {
	return &cli.Command{
		Name:  "run",
		Usage: "Start the server",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "config",
				Aliases: []string{"c"},
				Usage:   "Path to the config file",
				Value:   "ikou.config.json",
			},
		},
		Action: func(c *cli.Context) error {
			utils.InitLogger("prod")
			defer utils.Logger.Sync()
			utils.ExtractConfigDetails(c.String("config"))

			if err := react.BuildCSS(); err != nil {
				utils.Logger.Sugar().Fatalf("Failed to build CSS: %v", err)
				return err
			}
			app.StartServer(false)
			return nil
		},
	}
}
