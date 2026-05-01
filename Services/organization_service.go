package Services

import (
	"quay-go-api/Common"
	"quay-go-api/Common/Errors"
	"quay-go-api/Entities/Dto"
	"quay-go-api/Entities/Models"
	"quay-go-api/Repositories"
	"quay-go-api/Services/Auth"
	"quay-go-api/Services/Avatar"
)

func GetUserOrganizations(userId int, userScopes []Auth.Scope) ([]Dto.UserOrganization, error) {
	// Get the current user
	currentUser, err := Repositories.SelectUserById(userId)
	if err != nil {
		switch err.Error() {
		case "record not found":
			return nil, Errors.CurrentUserNotFound()
		default:
			return nil, err
		}
	}

	orgsModel, err := Repositories.GetUserOrganizations(userId)
	if err != nil {
		switch err.Error() {
		case "record not found": // if no result, is not an error, just return empty list
			return []Dto.UserOrganization{}, nil
		default:
			return nil, err
		}
	}

	// Convert userModel into UserOrganization Dto
	organizations := Common.ConvertUserModelsToDto(orgsModel, currentUser, userScopes)
	return organizations, nil
}

func GetOrganizationDetailsByName(orgName string, userId int, userScopes []Auth.Scope) (Dto.Organization, error) {
	// Get the current user
	currentUser, err := Repositories.SelectUserById(userId)
	if err != nil {
		switch err.Error() {
		case "record not found":
			return Dto.Organization{}, Errors.CurrentUserNotFound()
		default:
			return Dto.Organization{}, err
		}
	}

	// Get the organization with details
	orgModel, err := Repositories.GetOrganizationDetailsByName(orgName)
	if err != nil {
		switch err.Error() {
		case "record not found":
			return Dto.Organization{}, Errors.OrganizationNotFound(orgName)
		default:
			return Dto.Organization{}, err
		}
	}

	var teamsDto []Dto.Team

	// check if the current user is member of the organization (is in a team of the organization) and get his team role
	var userIsOrgAdmin bool = false
	var userIsOrgMember bool = false

	// If the user has scope super:user
	if Common.HasScope(userScopes, Auth.SuperUser) {
		userIsOrgAdmin = true
	}

	for _, orgTeam := range orgModel.Teams {
		for _, teamMember := range orgTeam.Members {
			if teamMember.User.ID == currentUser.ID {
				userIsOrgMember = true

				// Check if the user's team is the onwer team (team with role 'owner')
				if orgTeam.Role.Name == "owner" {
					userIsOrgAdmin = true
				}
			}
		}

		teamsDto = append(teamsDto, Dto.Team{
			Name:         orgTeam.Name,
			Description:  orgTeam.Description,
			Role:         orgTeam.Role.Name,
			Avatar:       Avatar.GetAvatarForTeam(orgTeam),
			CanView:      canViewTeams(userId, orgTeam, userScopes),
			MembersCount: len(orgTeam.Members),
			IsSynced:     false, // TODO: get if the team is synced
		})
	}

	orgDetailDto := Dto.Organization{
		Name:                orgModel.Username,
		Avatar:              Avatar.GetAvatarForOrg(orgModel),
		IsAdmin:             userIsOrgAdmin,
		IsMember:            userIsOrgMember,
		Teams:               teamsDto,
		InvoiceEmail:        orgModel.InvoiceEmail,
		InvoiceEmailAddress: Dto.NullString(orgModel.InvoiceEmailAddress),
		TagExpirationS:      orgModel.RemovedTagExpirationS,
		IsFreeAccount:       !orgModel.StripeId.Valid || orgModel.StripeId.String == "",
	}

	return orgDetailDto, nil
}

/*
canViewTeams checks if the user can view the team
A user can view a team if:
1. They are a member of that team (any role)
2. They are the scope org:admin
*/
func canViewTeams(userId int, team Models.Team, userScopes []Auth.Scope) bool {
	if team.Members == nil {
		panic("team members should be preloaded")
	}

	if Auth.Can(Auth.OrgAdmin, userScopes) {
		return true
	}
	for _, teamMember := range team.Members {
		if teamMember.UserId == userId {
			return true
		}
	}
	return false
}
