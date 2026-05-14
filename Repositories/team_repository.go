package Repositories

import (
	"gorm.io/gorm"
	"quay-go-api/Database"
	"quay-go-api/Entities/Models"
)

func CreateTeam(team Models.Team) (Models.Team, error) {
	err := Database.DB.Create(&team).Error
	if err != nil {
		return Models.Team{}, err
	}

	return team, nil
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
