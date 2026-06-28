package Repositories

import (
	"quay-go-api/Database"
	"quay-go-api/Entities/Models"
)

func GetRepositoryUserPermission(repositoryId int, userId int) (Models.RepositoryPermission, error) {
	var permission Models.RepositoryPermission
	err := Database.DB.Preload("Role").Where("repository_id = ? AND user_id = ?", repositoryId, userId).First(&permission).Error
	return permission, err
}
