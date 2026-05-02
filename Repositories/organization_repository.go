package Repositories

import (
	"gorm.io/gorm"
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

func GetOrganizationDetailsById(orgId int) (Models.User, error) {
	organization := Models.User{}

	err := Database.DB.
		Preload("Teams").
		Preload("Teams.Role").
		Preload("Teams.Members").
		Preload("Teams.Members.User").
		Where("organization = ?", true).
		First(&organization, orgId).
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

	err := Database.DB.
		Table("user AS organization_user").
		Distinct("organization_user.*").
		Joins("JOIN team ON team.organization_id = organization_user.id").
		Joins("JOIN teammember ON teammember.team_id = team.id").
		Joins("JOIN user AS member_user ON member_user.id = teammember.user_id").
		Where("organization_user.organization = ?", true).
		Where("member_user.id = ?", userId).
		Find(&organizations).Error

	return organizations, err
}

func CreateOrganization(organization Models.User) (Models.User, error) {
	err := Database.DB.Create(&organization).Error
	if err != nil {
		return Models.User{}, err
	}

	return organization, nil
}

func CreateOrganizationWithOwnerTeamTransaction(organization Models.User, userId int) (Models.User, error) {
	err := Database.DB.Transaction(func(tx *gorm.DB) error {
		// 1. Create the organization
		if err := tx.Create(&organization).Error; err != nil {
			return err // rollback
		}

		// 2. Define team model
		teamModel := Models.Team{
			Name:           "owners",
			OrganizationId: organization.ID,
			RoleId:         1, // 1 = admin, 2 = creator, 3 = member
		}

		// 3. Create the owner team
		if err := tx.Create(&teamModel).Error; err != nil {
			return err // rollback
		}

		// 4. Define team member model
		teamMemberModel := Models.TeamMember{
			UserId: userId,
			TeamId: teamModel.ID,
		}

		// 5. Add the user in the team
		if err := tx.Create(&teamMemberModel).Error; err != nil {
			return err
		}

		return nil // commit
	})

	if err != nil {
		return Models.User{}, err
	}
	return organization, nil
}

func DeleteOrganizationTransaction(orgId int) error {
	err := Database.DB.Transaction(func(tx *gorm.DB) error {
		// 1. Get Organization's teams id
		teamIdsSubQuery := tx.Model(&Models.Team{}).
			Select("id").
			Where("organization_id = ?", orgId)

		// 2. Remove all team members linked to the organization's teams
		if err := tx.Where("team_id IN (?)", teamIdsSubQuery).Delete(&Models.TeamMember{}).Error; err != nil {
			return err
		}

		// 3. Remove all teams linked to the organization
		if err := tx.Where("organization_id = ?", orgId).Delete(&Models.Team{}).Error; err != nil {
			return err
		}

		// 4. Remove the organization itself
		deleteResult := tx.Where("id = ?", orgId).
			Where("organization = ?", true).
			Delete(&Models.User{})

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

func UpdateOrganizationFieldsById(orgId int, fields map[string]interface{}) error {
	if len(fields) == 0 {
		return nil
	}

	result := Database.DB.Model(&Models.User{}).
		Where("id = ?", orgId).
		Where("organization = ?", true).
		Updates(fields)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}
