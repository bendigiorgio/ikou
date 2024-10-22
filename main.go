package main

import (
	"github.com/bendigiorgio/ikou/internal/app"
	"github.com/bendigiorgio/ikou/internal/app/utils"
)

func main() {
	utils.ExtractConfigDetails("ikou.config.json")
	app.StartServer()
}
