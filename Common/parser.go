package Common

import (
	"quay-go-api/Entities/Dto"
	"quay-go-api/Entities/Models"
	"quay-go-api/Services/Auth"
	"quay-go-api/Services/Avatar"
	"strings"
)

func ConvertScopeStringInAuthScopes(scopesStr string) []Auth.Scope {
	scopeIDs := strings.Split(scopesStr, " ")
	var scopes []Auth.Scope

	for _, scopeID := range scopeIDs {
		scopes = append(scopes, Auth.GetScopeFromID(scopeID))
	}
	return scopes
}

func ConvertUserModelsToDto(orgsModel []Models.User, currentUser Models.User, userScopes []Auth.Scope) []Dto.UserOrganization {
	var orgs []Dto.UserOrganization

	for _, org := range orgsModel {
		orgs = append(orgs, Dto.UserOrganization{
			Name:               org.Username,
			Avatar:             Avatar.GetAvatarForOrg(org),
			CanCreateRepo:      Auth.Can(Auth.CreateRepo, userScopes),
			Public:             false, // TODO: check if the org name not in list of public Namespaces
			IsOrgAdmin:         Auth.Can(Auth.OrgAdmin, userScopes),
			PreferredNamespace: !(!currentUser.StripeId.Valid || currentUser.StripeId.String == ""),
		})
	}

	return orgs
}

func ConvertUserModelToOrganizationDto(orgDetailsModel Models.User, currentUserId int, userScopes []Auth.Scope) Dto.Organization {
	var teamsDto []Dto.Team

	// check if the current user is member of the organization (is in a team of the organization) and get his team role
	var userIsOrgAdmin bool = false
	var userIsOrgMember bool = false

	// If the user has scope super:user
	if HasScope(userScopes, Auth.SuperUser) {
		userIsOrgAdmin = true
	}

	for _, orgTeam := range orgDetailsModel.Teams {
		for _, teamMember := range orgTeam.Members {
			if teamMember.User.ID == currentUserId {
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
			CanView:      canViewTeams(currentUserId, orgTeam, userScopes),
			MembersCount: len(orgTeam.Members),
			IsSynced:     false, // TODO: get if the team is synced
		})
	}

	orgDetailDto := Dto.Organization{
		Name:                orgDetailsModel.Username,
		Avatar:              Avatar.GetAvatarForOrg(orgDetailsModel),
		IsAdmin:             userIsOrgAdmin,
		IsMember:            userIsOrgMember,
		Teams:               teamsDto,
		InvoiceEmail:        orgDetailsModel.InvoiceEmail,
		InvoiceEmailAddress: Dto.NullString(orgDetailsModel.InvoiceEmailAddress),
		TagExpirationS:      orgDetailsModel.RemovedTagExpirationS,
		IsFreeAccount:       !orgDetailsModel.StripeId.Valid || orgDetailsModel.StripeId.String == "",
	}

	return orgDetailDto
}
