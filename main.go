package main

import (
	"os"
	"quay-go-api/api"
	"quay-go-api/service/logger"
)

func main() {
	if logLevel := os.Getenv("LOG_LEVEL"); logLevel != "" {
		logger.SetLevel(logger.StringToLevel(logLevel))
	} else {
		logger.SetLevel(logger.LevelDebug)
	}

	api.StartServer()
}
