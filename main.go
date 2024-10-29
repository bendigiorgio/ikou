package main

import (
	"github.com/bendigiorgio/ikou/internal/app"
	"github.com/bendigiorgio/ikou/internal/app/react"
	"github.com/bendigiorgio/ikou/internal/app/utils"
)

func main() {
	defer utils.Logger.Sync()
	utils.ExtractConfigDetails("ikou.config.json")

	devMode := true

	if devMode {
		go utils.WatchForConfigChanges("ikou.config.json")
	}

	if err := react.BuildCSS(); err != nil {
		utils.Logger.Sugar().Fatalf("Failed to build CSS: %v", err)
	}

	app.StartServer(devMode)
}
