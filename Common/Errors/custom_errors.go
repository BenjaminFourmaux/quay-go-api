package Errors

import (
	"net/http"
	"quay-go-api/Services/Auth"
)

func Unauthorized() *ApiError {
	return &ApiError{
		StatusCode: http.StatusUnauthorized,
		Err: ErrorResponse{
			Error: ErrorDetails{
				Code:    "unauthorized",
				Message: "You cannot perform this action because your haven't provided a authentication token",
			},
		},
	}
}

func ForbiddenInvalidToken() *ApiError {
	return &ApiError{
		StatusCode: http.StatusForbidden,
		Err: ErrorResponse{
			Error: ErrorDetails{
				Code:    "invalid_token",
				Message: "The provided token is invalid or has expired",
			},
		},
	}
}

func ForbiddenNoRequiredScope(scopes []Auth.Scope) *ApiError {
	missingScopes := ""
	for i, scope := range scopes {
		missingScopes += scope.Name
		if i < len(scopes)-1 {
			missingScopes += ", "
		}
	}

	return &ApiError{
		StatusCode: http.StatusForbidden,
		Err: ErrorResponse{
			Error: ErrorDetails{
				Code:    "insufficient_scope",
				Message: "you do not have the required permissions (" + missingScopes + ")to access this resource",
			},
		},
	}
}
