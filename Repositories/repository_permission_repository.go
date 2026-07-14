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
ListRepositoryPermissions Return list of user or team repository permissions by repository ID.
Parameter 'kind' can be empty for get both user and team permissions
*/
func ListRepositoryPermissions(repositoryId int, kind string) ([]Models.RepositoryPermission, error) {
	var permissions []Models.RepositoryPermission

	// Not beautiful, rip my Common.InlineIf()
	var whereKind string
	if kind == "team" {
		whereKind = "team_id IS NOT NULL"
	} else if kind == "user" {
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

func GetUserRepositoryPermissionByUsername(repositoryId int, username string) (Models.RepositoryPermission, error) {
	var permission Models.RepositoryPermission
	err := Database.DB.
		Preload("Role").
		InnerJoins("User", Database.DB.Where(&Models.User{Username: username})).
		Where("repositorypermission.repository_id = ?", repositoryId).
		First(&permission).
		Error
	return permission, err
}

func GetTeamRepositoryPermissionByTeamname(repositoryId int, teamname string) (Models.RepositoryPermission, error) {
	var permission Models.RepositoryPermission
	err := Database.DB.
		Preload("Role").
		InnerJoins("Team", Database.DB.Where(&Models.Team{Name: teamname})).
		Where("repositorypermission.repository_id = ?", repositoryId).
		First(&permission).
		Error
	return permission, err
}

func UpdateRepositoryPermission(permission Models.RepositoryPermission) error {
	err := Database.DB.Model(&permission).Updates(&permission).Error
	return err
}

func DeleteRepositoryPermission(permission Models.RepositoryPermission) error {
	err := Database.DB.Delete(&permission).Error
	return err
}
