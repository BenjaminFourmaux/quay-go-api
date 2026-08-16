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

// @externalDocs.description Wiki
// @externalDocs.url https://github.com/BenjaminFourmaux/quay-go-api/wiki

// @tag.name Messages
// @tag.description global messages
// @tag.name Manifest
// @tag.description Manage the manifest of a repository
// @tag.name Organization
// @tag.description Manage organizations
// @tag.name Members
// @tag.description Manage organization's members
// @tag.name Permission
// @tag.description Manage repository permissions
// @tag.name Prototype
// @tag.description Manage default permissions added to repositories
// @tag.name Repository
// @tag.description List, create and manage repositories
// @tag.name Robot
// @tag.description Manage user and organization robot accounts
// @tag.name SecurityScan
// @tag.description List and manage repository vulnerabilities and other security information
// @tag.name Tag
// @tag.description Manage the tags of a repository
// @tag.name Team
// @tag.description Create, list and manage and organization's teams
// @tag.name Users
// @tag.description Manage users
// @tag.name Avatar
// @tag.description Get user, team or organization avatar

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
