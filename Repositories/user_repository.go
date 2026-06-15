package Repositories

import (
	"quay-go-api/Database"
	"quay-go-api/Entities/Models"
)

func GetUserById(userId int) (Models.User, error) {
	var user Models.User

	err := Database.DB.First(&user, userId).Error

	return user, err
}

/*
GetUserByName Get a user (not an org) from username
*/
func GetUserByName(userName string) (Models.User, error) {
	var user Models.User

	err := Database.DB.
		Where("organization = ?", false).
		Where("username = ?", userName).
		First(&user).
		Error

	return user, err
}

/*
GetUserByIdWithUserInformation Get a user by id and with federated logins, associate login service and user Prompts
*/
func GetUserByIdWithUserInformation(userId int) (Models.User, error) {
	var user Models.User

	err := Database.DB.
		Preload("FederatedLogins.Service").
		Preload("Prompts.Kind").
		First(&user, userId).
		Error

	return user, err
}

/*
GetUserOrOrganizationByName Get both User or Organization by his name
*/
func GetUserOrOrganizationByName(organizationName string) (Models.User, error) {
	var user Models.User

	err := Database.DB.Where("username = ?", organizationName).First(&user).Error

	return user, err
}
