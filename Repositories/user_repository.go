package Repositories

import (
	"quay-go-api/Database"
	"quay-go-api/Entities/Models"
)

func SelectUserById(userId int) Models.User {
	var user Models.User
	Database.DB.First(&user, userId)
	return user
}

/*
GetUserByIdWithUserInformation Get a user by id and with federated logins, associate login service and user Prompts
*/
func GetUserByIdWithUserInformation(userId int) Models.User {
	var user Models.User
	Database.DB.
		Preload("FederatedLogins.Service").
		Preload("Prompts.Kind").
		First(&user, userId)
	return user
}
