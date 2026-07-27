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

func CreatePermissionPrototype(prototype *Models.PermissionPrototype) (*Models.PermissionPrototype, error) {
	err := Database.DB.Create(prototype).Error
	if err != nil {
		return nil, err
	}
	err = Database.DB.Preload("Organization").
		Preload("ActivatingUser").
		Preload("DelegateUser").
		Preload("DelegateTeam").
		Preload("Role").
		First(prototype, prototype.ID).Error
	return prototype, err
}

func DeletePermissionPrototype(prototypeId int) error {
	err := Database.DB.Delete(&Models.PermissionPrototype{}, prototypeId).Error
	return err
}
