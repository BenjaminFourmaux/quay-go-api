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

func MessageInvalidSeverity(wrongSeverity string) *ApiError {
	return &ApiError{
		StatusCode: http.StatusBadRequest,
		Err: ErrorResponse{
			Error: ErrorDetails{
				Code:    "message_invalid_severity",
				Message: "The provided severity '" + wrongSeverity + "' is invalid. Valid severity levels are: 'info', 'warning', 'error'",
			},
		},
	}
}
