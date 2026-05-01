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

func ConvertUserModelsToDto(orgsModel []Models.User, currentUser Models.User, userScopes []Auth.Scope) []Dto.Organization {
	var orgs []Dto.Organization

	for _, org := range orgsModel {
		orgs = append(orgs, Dto.Organization{
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
