package Common

import (
	"quay-go-api/Entities/Models"
	"quay-go-api/Services/Auth"
)

func HasScope(scopes []Auth.Scope, scope Auth.Scope) bool {
	for _, scp := range scopes {
		if scp.ID == scope.ID {
			return true
		}
	}
	return false
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
