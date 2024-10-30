package utils

import (
	"log"
	"os"

	"go.uber.org/zap"
)

var Logger *zap.Logger
var internal_mode string

func InitLogger(mode string) {
	if Logger != nil {
		return
	}
	internal_mode = mode
	logPath := "storage/logs/ikou.log"

	if _, err := os.Stat(logPath); os.IsNotExist(err) {
		if _, err := os.Create(logPath); err != nil {
			log.Fatalf("failed to create log file: %v", err)
		}
	}

	config := zap.NewDevelopmentConfig()
	if mode == "prod" {
		config = zap.NewProductionConfig()
	}
	config.OutputPaths = []string{"stdout", logPath}
	config.ErrorOutputPaths = []string{"stderr", logPath}

	logger, err := config.Build()
	if err != nil {
		log.Fatalf("failed to initialize zap logger: %v", err)
		panic("failed to initialize zap logger")
	}

	Logger = logger
}

func UpdateLogPath(newPath string) {
	if _, err := os.Stat(newPath); os.IsNotExist(err) {
		if _, err := os.Create(newPath); err != nil {
			log.Fatalf("failed to create log file: %v", err)
		}
	}

	config := zap.NewDevelopmentConfig()
	if internal_mode == "prod" {
		config = zap.NewProductionConfig()
	}

	config.OutputPaths = []string{"stdout", newPath}
	config.ErrorOutputPaths = []string{"stderr", newPath}

	newLogger, err := config.Build()
	if err != nil {
		Logger.Sugar().Errorf("failed to update log path: %v", err)
		return
	}
	Logger = newLogger
}
