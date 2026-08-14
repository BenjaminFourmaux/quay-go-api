package Repositories

import (
	"fmt"
	"gorm.io/gorm"
	"quay-go-api/Database"
	"quay-go-api/Entities/Models"
)

/*
GetUserOrOrgRobots retrieves all robot accounts associated with a specific user or organization (an org is a user in bd). It allows for optional inclusion of robot tokens and repository permissions.
*/
func GetUserOrOrgRobots(userID int, includeToken bool, includeRepositories bool) ([]Models.User, error) {
	var robots []Models.User

	user, err := GetUserById(userID)
	if err != nil {
		return nil, err
	}

	usernamePattern := user.Username + "+%"
	query := Database.DB.
		Preload("RobotAccountMetadata").
		Where("robot = ?", true).
		Where("username LIKE ?", usernamePattern)

	if includeToken {
		query = query.Preload("RobotAccountToken")
	}

	if includeRepositories {
		query = query.
			Preload("RepositoryPermissions", "repository_id IS NOT NULL").
			Preload("RepositoryPermissions.Repository")
	}

	err = query.Find(&robots).Error

	return robots, err
}

func GetRobotByName(robotName string, userId int) (Models.User, error) {
	var robot Models.User

	err := Database.DB.
		Where("username = ?", robotName). // TODO: maybe find a method to identify robot accounts form a user/org by id
		First(&robot).Error

	return robot, err
}

func CreateRobotAccount(robotToCreate Models.User, robotMetadata Models.RobotAccountMetadata, robotToken Models.RobotAccountToken, federatedLogin Models.FederatedLogin) (Models.User, error) {
	// Start a transaction
	err := Database.DB.Transaction(func(tx *gorm.DB) error {
		// 1. Create the robot user
		if err := tx.Create(&robotToCreate).Error; err != nil {
			return err
		}

		// 2. Set user id to robotMetadata and token
		robotMetadata.RobotAccountID = robotToCreate.ID
		robotToken.RobotAccountID = robotToCreate.ID
		federatedLogin.UserId = robotToCreate.ID
		federatedLogin.ServiceIdent = fmt.Sprintf("robot:%d", robotToCreate.ID)

		// 3. Create the robot metadata
		if err := tx.Create(&robotMetadata).Error; err != nil {
			return err
		}

		// 4. Create the robot token
		if err := tx.Create(&robotToken).Error; err != nil {
			return err
		}

		// 5. Create the federated login
		if err := tx.Create(&federatedLogin).Error; err != nil {
			return err
		}

		return nil // commit
	})

	// Assemble metadata and token into robotToCreate
	robotToCreate.RobotAccountMetadata = &robotMetadata
	robotToCreate.RobotAccountToken = &robotToken
	
	return robotToCreate, err
}
