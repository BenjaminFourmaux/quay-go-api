package Repositories

import (
	"quay-go-api/Database"
	"quay-go-api/Entities/Models"
)

func GetUserRobots(userID int, includeToken bool, includeRepositories bool) ([]Models.User, error) {
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
