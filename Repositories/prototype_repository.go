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

func GetOrganizationPrototypeByUUID(orgId int, uuid string) (Models.PermissionPrototype, error) {
	var prototype Models.PermissionPrototype
	err := Database.DB.
		Where("org_id = ? AND uuid = ?", orgId, uuid).
		First(&prototype).
		Error
	return prototype, err
}

func GetOrganizationPrototypeDetailsByUUID(orgId int, uuid string) (Models.PermissionPrototype, error) {
	var prototype Models.PermissionPrototype
	err := Database.DB.
		Preload("Organization").
		Preload("ActivatingUser").
		Preload("DelegateUser").
		Preload("DelegateTeam").
		Preload("Role").
		Where("org_id = ? AND uuid = ?", orgId, uuid).
		First(&prototype).
		Error
	return prototype, err
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

func UpdatePermissionPrototype(prototype *Models.PermissionPrototype) (*Models.PermissionPrototype, error) {
	err := Database.DB.Save(prototype).Error
	if err != nil {
		return nil, err
	}
	return prototype, nil
}

func DeletePermissionPrototype(prototypeId int) error {
	err := Database.DB.Delete(&Models.PermissionPrototype{}, prototypeId).Error
	return err
}

func CheckIfPermissionPrototypeExists(orgId int, activatingUserId *int, delegateUserId *int, delegateTeamId *int) (Models.PermissionPrototype, error) {
	var prototype Models.PermissionPrototype
	query := Database.DB.
		Where("org_id = ?", orgId)

	if activatingUserId != nil {
		query = query.Where("activating_user_id = ?", *activatingUserId)
	}
	if delegateUserId != nil {
		query = query.Where("delegate_user_id = ?", *delegateUserId)
	}
	if delegateTeamId != nil {
		query = query.Where("delegate_team_id = ?", *delegateTeamId)
	}

	err := query.First(&prototype).Error
	return prototype, err
}
