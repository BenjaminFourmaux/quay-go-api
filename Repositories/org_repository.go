package Repositories

import (
	"quay-go-api/Database"
	"quay-go-api/Entities/Models"
)

/*
GetUserOrganizations returns the organization users accessible to the specified username
through team membership.
*/
func GetUserOrganizations(username string) ([]Models.User, error) {
	organizations := []Models.User{}

	// "user" is a reserved keyword in both PostgreSQL (CURRENT_USER) and MySQL.
	// It must be quoted to refer to the actual table.
	db := Database.DB
	var userTable string
	if db.Dialector.Name() == "postgres" {
		userTable = `"user"`
	} else {
		userTable = "`user`"
	}

	err := db.
		Table(userTable+" AS organization_user").
		Distinct("organization_user.*").
		Joins("JOIN team ON team.organization_id = organization_user.id").
		Joins("JOIN teammember ON teammember.team_id = team.id").
		Joins("JOIN "+userTable+" AS member_user ON member_user.id = teammember.user_id").
		Where("organization_user.organization = ?", true).
		Where("member_user.username = ?", username).
		Find(&organizations).Error

	return organizations, err
}
