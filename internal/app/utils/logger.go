package utils

import "go.uber.org/zap"

var Logger *zap.Logger

func init() {
	config := zap.NewDevelopmentConfig()
	config.OutputPaths = []string{"stdout", "ikou.log"}
	logger, _ := config.Build()
	Logger = logger
}
