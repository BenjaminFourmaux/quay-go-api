package Common

import (
	"database/sql"
	"quay-go-api/Entities/Dto"
	"quay-go-api/Entities/Models"
	"quay-go-api/Services/Auth"
	"quay-go-api/Services/Avatar"
	"strings"
	"time"
)

func ParseTime(datetime string) time.Time {
	t, _ := time.Parse(time.RFC3339, datetime)
	return t
}

func ConvertSQLNullTimeToTime(nullTime sql.NullTime) *time.Time {
	if nullTime.Valid {
		return &nullTime.Time
	}
	return nil
}

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

		teamsDto = append(teamsDto, ConvertTeamModelToDto(orgTeam, currentUserId, userScopes))
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

func ConvertTeamModelToDto(teamModel Models.Team, userId int, userScopes []Auth.Scope) Dto.Team {
	return Dto.Team{
		Name:         teamModel.Name,
		Description:  teamModel.Description,
		Role:         teamModel.Role.Name,
		Avatar:       Avatar.GetAvatarForTeam(teamModel),
		CanView:      CanViewTeams(userId, teamModel, userScopes),
		MembersCount: len(teamModel.Members),
		IsSynced:     false, // TODO: get if the team is synced
	}
}

func ConvertRepositoryPermissionModelToDto(repositoryPermissionModel Models.RepositoryPermission, kind string) Dto.RepositoryPermission {
	permission := Dto.RepositoryPermission{
		Role: repositoryPermissionModel.Role.Name,
	}

	if kind == "user" {
		permission.Name = repositoryPermissionModel.User.Username
		permission.Avatar = Avatar.GetAvatarForUser(*repositoryPermissionModel.User)
		permission.IsRobot = &repositoryPermissionModel.User.Robot
	} else if kind == "team" {
		permission.Name = repositoryPermissionModel.Team.Name
		permission.Avatar = Avatar.GetAvatarForTeam(*repositoryPermissionModel.Team)
	}
	return permission
}

func ConvertPermissionPrototypeModelToDto(prototypeModel Models.PermissionPrototype, activatingUserOrgMember bool, delegateOrgMember bool) Dto.Prototype {
	prototype := Dto.Prototype{
		Id:   prototypeModel.UUID,
		Role: prototypeModel.Role.Name,
	}

	if prototypeModel.ActivatingUserId != nil && prototypeModel.ActivatingUser != nil { // Always true
		prototype.ActivatingUser = Dto.ActivatingUser{
			Name:        prototypeModel.ActivatingUser.Username,
			IsRobot:     prototypeModel.ActivatingUser.Robot,
			Kind:        "user",
			IsOrgMember: activatingUserOrgMember,
			Avatar:      Avatar.GetAvatarForUser(*prototypeModel.ActivatingUser),
		}
	}

	// Delegate to a user
	if prototypeModel.DelegateUserId != nil && prototypeModel.DelegateUser != nil {
		prototype.Delegate = Dto.Delegate{
			Name:        prototypeModel.DelegateUser.Username,
			IsRobot:     prototypeModel.DelegateUser.Robot,
			Kind:        "user",
			IsOrgMember: delegateOrgMember,
			Avatar:      Avatar.GetAvatarForUser(*prototypeModel.DelegateUser),
		}
	}

	// Delegate to a team
	if prototypeModel.DelegateTeamId != nil && prototypeModel.DelegateTeam != nil {
		prototype.Delegate = Dto.Delegate{
			Name:        prototypeModel.DelegateTeam.Name,
			IsRobot:     false, // A team cannot be a robot account
			Kind:        "team",
			IsOrgMember: delegateOrgMember,
			Avatar:      Avatar.GetAvatarForTeam(*prototypeModel.DelegateTeam),
		}
	}

	return prototype
}
