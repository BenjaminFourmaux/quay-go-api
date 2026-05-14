package Repositories

import (
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
