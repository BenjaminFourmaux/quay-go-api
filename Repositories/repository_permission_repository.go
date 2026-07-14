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

/*
ListRepositoryPermissions Return list of user or team repository permissions by repository ID
*/
func ListRepositoryPermissions(repositoryId int, kind string) ([]Models.RepositoryPermission, error) {
	var permissions []Models.RepositoryPermission

	// Not beautiful, rip my Common.InlineIf()
	var whereKind string
	if kind == "team" {
		whereKind = "team_id IS NOT NULL"
	} else {
		whereKind = "user_id IS NOT NULL"
	}

	err := Database.DB.
		Preload("Role").
		Preload("User").
		Preload("Team").
		Where("repository_id = ?", repositoryId).
		Where(whereKind).
		Find(&permissions).
		Error
	return permissions, err
}
