package Repositories

import (
	"quay-go-api/Database"
	"quay-go-api/Entities/Models"
)

func GetOrganizationByName(orgName string) (Models.User, error) {
	organization := Models.User{}

	err := Database.DB.
		Where("organization = ?", true).
		Where("username = ?", orgName).
		First(&organization).
		Error

	return organization, err
}

func GetOrganizationDetailsByName(orgName string) (Models.User, error) {
	organization := Models.User{}

	err := Database.DB.
		Preload("Teams").
		Preload("Teams.Role").
		Preload("Teams.Members").
		Preload("Teams.Members.User").
		Where("organization = ?", true).
		Where("username = ?", orgName).
		First(&organization).
		Error

	return organization, err
}

/*
GetUserOrganizations returns the organization users accessible to the specified username
through team membership.
*/
func GetUserOrganizations(userId int) ([]Models.User, error) {
	organizations := []Models.User{} // Orgs are in table user

	// "user" is a reserved keyword in both PostgreSQL (CURRENT_USER) and MySQL.
	// It must be quoted to refer to the actual table.
	db := Database.DB

	// Get the user table name depending on the connected database type
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
		Where("member_user.id = ?", userId).
		Find(&organizations).Error

	return organizations, err
}
