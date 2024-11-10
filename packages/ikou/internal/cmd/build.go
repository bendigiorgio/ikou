package cmd

import (
	"github.com/bendigiorgio/ikou/internal/app/ssg"
	"github.com/bendigiorgio/ikou/internal/app/utils"
	"github.com/urfave/cli/v2"
)

func GetBuildCommand() *cli.Command {
	return &cli.Command{
		Name:  "build",
		Usage: "Build Static Files",
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
			ssg.GenerateStaticSite()
			utils.Logger.Sugar().Info("Static site generated. Run the command `ikou serve` to serve the static site.")
			return nil
		},
	}
}
