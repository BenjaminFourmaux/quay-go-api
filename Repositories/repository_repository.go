package Repositories

import (
	"quay-go-api/Database"
	"quay-go-api/Entities/Models"
)

func SelectRepositories(userId int, organizationOrUserId *int) ([]Models.Repository, error) {
	var repositories []Models.Repository
	query := Database.DB.
		Preload("NamespaceUser").
		Preload("Visibility").
		Preload("Kind").
		Preload("Stars", "user_id = ?", userId)

	if organizationOrUserId != nil {
		query = query.Where("namespace_user_id = ?", *organizationOrUserId)
	}

	err := query.Find(&repositories).Error
	return repositories, err
}

func GetOrganizationRepositoriesByOrgId(organizationId int, userId int) ([]Models.Repository, error) {
	var repositories []Models.Repository
	err := Database.DB.
		Preload("Visibility").
		Preload("Kind").
		Preload("Stars", "user_id = ?", userId).
		Where("namespace_user_id = ?", organizationId).
		Find(&repositories).
		Error
	return repositories, err
}
