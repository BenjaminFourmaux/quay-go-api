package Repositories

import (
	"quay-go-api/Database"
	"quay-go-api/Entities/Models"
)

func GetOrganizationPrototypes(orgId int) ([]Models.PermissionPrototype, error) {
	var prototypes []Models.PermissionPrototype
	err := Database.DB.
		Preload("Organization").
		Preload("ActivatingUser").
		Preload("DelegateUser").
		Preload("DelegateTeam").
		Preload("Role").
		Where("org_id = ?", orgId).
		Find(&prototypes).
		Error
	return prototypes, err
}
