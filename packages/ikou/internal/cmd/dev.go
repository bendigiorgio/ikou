package cmd

import (
	"github.com/bendigiorgio/ikou/internal/app"
	"github.com/bendigiorgio/ikou/internal/app/react"
	"github.com/bendigiorgio/ikou/internal/app/utils"
	"github.com/urfave/cli/v2"
)

func GetDevCommand() *cli.Command {
	return &cli.Command{
		Name:  "dev",
		Usage: "Start the development server",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "config",
				Aliases: []string{"c"},
				Usage:   "Path to the config file",
				Value:   "ikou.config.json",
			},
		},
		Action: func(c *cli.Context) error {
			utils.InitLogger("dev")
			defer utils.Logger.Sync()
			utils.ExtractConfigDetails(c.String("config"))
			go utils.WatchForConfigChanges(c.String("config"))

			if err := react.BuildCSS(); err != nil {
				utils.Logger.Sugar().Fatalf("Failed to build CSS: %v", err)
				return err
			}

			app.StartServer(true)
			return nil
		},
	}
}
