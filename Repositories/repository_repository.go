package Repositories

import (
	"quay-go-api/Database"
	"quay-go-api/Entities/Models"
)

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
