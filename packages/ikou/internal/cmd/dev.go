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
		Action: func(c *cli.Context) error {
			go utils.WatchForConfigChanges("ikou.config.json")

			if err := react.BuildCSS(); err != nil {
				utils.Logger.Sugar().Fatalf("Failed to build CSS: %v", err)
				return err
			}

			app.StartServer(true)
			return nil
		},
	}
}
