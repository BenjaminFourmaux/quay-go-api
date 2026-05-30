package Repositories

import (
	"quay-go-api/Database"
	"quay-go-api/Entities/Models"
)

func CreateTeamMember(teamMember Models.TeamMember) (Models.TeamMember, error) {
	err := Database.DB.Create(&teamMember).Error
	if err != nil {
		return Models.TeamMember{}, err
	}
	return teamMember, nil
}

func DeleteTeamMember(teamMember Models.TeamMember) error {
	return Database.DB.Delete(&teamMember).Error
}
