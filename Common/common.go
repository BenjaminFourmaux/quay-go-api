package Common

import "quay-go-api/Services/Auth"

func HasScope(scopes []Auth.Scope, scope Auth.Scope) bool {
	for _, scp := range scopes {
		if scp.ID == scope.ID {
			return true
		}
	}
	return false
}
