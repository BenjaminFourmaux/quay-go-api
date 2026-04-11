package main

import (
	"os"
	"quay-go-api/api"
	_ "quay-go-api/docs"
	"quay-go-api/service/logger"
)

// @title Quay Go API
// @version 1.0
// @description Quay registry API implemented in Go
func main() {
	if logLevel := os.Getenv("LOG_LEVEL"); logLevel != "" {
		logger.SetLevel(logger.StringToLevel(logLevel))
	} else {
		logger.SetLevel(logger.LevelDebug)
	}

	api.StartServer()
}
