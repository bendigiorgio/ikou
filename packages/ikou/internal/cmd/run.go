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
		Action: func(c *cli.Context) error {
			if err := react.BuildCSS(); err != nil {
				utils.Logger.Sugar().Fatalf("Failed to build CSS: %v", err)
				return err
			}
			app.StartServer(false)
			return nil
		},
	}
}
