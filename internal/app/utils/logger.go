package utils

import "go.uber.org/zap"

var Logger *zap.Logger

func init() {
	logPath := GlobalConfig.LogPath

	config := zap.NewDevelopmentConfig()
	config.OutputPaths = []string{"stdout", logPath}
	logger, _ := config.Build()
	Logger = logger
}
