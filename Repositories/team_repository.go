package Repositories

import (
	"gorm.io/gorm"
	"quay-go-api/Database"
	"quay-go-api/Entities/Models"
)

func GetTeamById(teamId int) (Models.Team, error) {
	var team Models.Team
	err := Database.DB.
		Preload("Role").
		Preload("Members").
		Preload("Members.User").
		First(&team, teamId).
		Error

	return team, err
}

func GetTeamDetailsById(teamId int) (Models.Team, error) {
	var team Models.Team
	err := Database.DB.
		Preload("Role").
		Preload("Members").
		Preload("Members.User").
		Preload("TeamMemberInvites").
		First(&team, teamId).
		Error
	return team, err
}

/*
GetOrganizationTeamsByOrgId returns organization teams with Role, Members, Members.User and TeamMemberInvites preloaded
*/
func GetOrganizationTeamsByOrgId(organizationId int) ([]Models.Team, error) {
	var teams []Models.Team
	err := Database.DB.
		Preload("Role").
		Preload("Members").
		Preload("Members.User").
		Preload("TeamMemberInvites").
		Where("organization_id = ?", organizationId).
		Find(&teams).
		Error
	return teams, err
}

func CreateTeam(team Models.Team) (Models.Team, error) {
	err := Database.DB.Create(&team).Error
	if err != nil {
		return Models.Team{}, err
	}

	return team, nil
}

func UpdateTeamFieldsById(teamId int, fields map[string]interface{}) error {
	if len(fields) == 0 {
		return nil
	}

	result := Database.DB.Model(&Models.Team{}).
		Where("id = ?", teamId).
		Updates(fields)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil

}

func DeleteTeamTransaction(teamId int) error {
	err := Database.DB.Transaction(func(tx *gorm.DB) error {
		// 1. Remove all team members linked to the team
		if err := tx.Where("team_id = ?", teamId).Delete(&Models.TeamMember{}).Error; err != nil {
			return err
		}

		// 2. Remove the team
		deleteResult := tx.Delete(&Models.Team{}, teamId)

		if deleteResult.Error != nil {
			return deleteResult.Error
		}

		if deleteResult.RowsAffected == 0 {
			return gorm.ErrRecordNotFound
		}

		return nil
	})
	return err
}
