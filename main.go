package main

import (
	"os"
	"quay-go-api/logger"
)

func main() {
	if logLevel := os.Getenv("LOG_LEVEL"); logLevel != "" {
		logger.SetLevel(logger.StringToLevel(logLevel))
	} else {
		logger.SetLevel(logger.LevelDebug)
	}

	logger.Debug("An error occurred")
}
