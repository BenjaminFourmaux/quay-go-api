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
CanViewTeams checks if the user can view the team
A user can view a team if:
1. They are a member of that team (any role)
2. They are the scope org:admin
*/
func CanViewTeams(userId int, team Models.Team, userScopes []Auth.Scope) bool {
	if Auth.Can(Auth.OrgAdmin, userScopes) {
		return true
	}

	if team.Members == nil {
		panic("team members should be preloaded")
	}
	for _, teamMember := range team.Members {
		if teamMember.UserId == userId {
			return true
		}
	}
	return false
}

func GetTeamRoleIdFromRoleName(roleName string) int {
	switch roleName {
	case "admin":
		return 1
	case "creator":
		return 2
	case "member":
		return 3
	default:
		return 0
	}
}

func MapRepositoryStateStr(stateId int) string {
	switch stateId {
	case 0:
		return "NORMAL" // Regular repo where all actions are possible
	case 1:
		return "READ_ONLY" // Only read actions, such as pull, are allowed regardless of specific user permissions
	case 2:
		return "MIRROR" // Equivalent to READ_ONLY except that mirror robot has write permission
	case 3:
		return "MARKED_FOR_DELETION" // Indicates the repository has been marked for deletion and should be hidden and unusable.
	case 4:
		return "ORG_MIRROR" // Equivalent to MIRROR but for repositories created via organization-level mirroring
	default:
		return "UNKNOWN"
	}
}

func GetRoleIdFromRoleName(roleName string) int {
	switch roleName {
	case "admin":
		return 1
	case "write":
		return 2
	case "read":
		return 3
	default:
		return 0
	}
}
