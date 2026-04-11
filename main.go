package main

import (
	"os"
	"quay-go-api/Api"
	"quay-go-api/Database"
	"quay-go-api/Services/Logger"
	_ "quay-go-api/docs"
)

// @title Quay Go API
// @version 1.0
// @description Quay registry API implemented in Go
// @securityDefinitions.apiKey ApiKeyAuth
// @in header
// @name Authorization
// @description API key authentication using the Authorization header. The value should be in the format "Bearer {token}"
func main() {
	if logLevel := os.Getenv("LOG_LEVEL"); logLevel != "" {
		Logger.SetLevel(Logger.StringToLevel(logLevel))
	} else {
		Logger.SetLevel(Logger.LevelDebug)
	}

	// Connect to the database
	Database.ConnectDatabase()

	// Start the HTTP server
	Api.StartServer()
}
