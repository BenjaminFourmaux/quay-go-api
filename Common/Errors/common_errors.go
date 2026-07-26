package Errors

import (
	"net/http"
	"quay-go-api/Services/Auth"
	"strings"
)

func InternalServerError() *ApiError {
	return &ApiError{
		StatusCode: http.StatusInternalServerError,
		Err: ErrorResponse{
			Error: ErrorDetails{
				Code:    "internal_server_error",
				Message: "An Internal Server Error was throw",
			},
		},
	}
}

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

func UnauthorizedInsufficientRole() *ApiError {
	return &ApiError{
		StatusCode: http.StatusForbidden,
		Err: ErrorResponse{
			Error: ErrorDetails{
				Code:    "insufficient_role",
				Message: "You cannot perform this action because your haven't correct role",
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

func BadRequest(msg string) *ApiError {
	return &ApiError{
		StatusCode: http.StatusBadRequest,
		Err: ErrorResponse{
			Error: ErrorDetails{
				Code:    "bad_request",
				Message: msg,
			},
		},
	}
}

func InvalidParameterValue(paramName string, allowedValues []string) *ApiError {
	quotedValues := make([]string, len(allowedValues))
	for i, val := range allowedValues {
		quotedValues[i] = "'" + val + "'"
	}

	return &ApiError{
		StatusCode: http.StatusBadRequest,
		Err: ErrorResponse{
			Error: ErrorDetails{
				Code:    "invalid_parameter_value",
				Message: "The provided parameter '" + paramName + "' has an invalid value. Allowed values: " + strings.Join(quotedValues, ", "),
			},
		},
	}
}

func CurrentUserNotFound() *ApiError {
	return &ApiError{
		StatusCode: http.StatusNotFound,
		Err: ErrorResponse{
			Error: ErrorDetails{
				Code:    "current_user_not_found",
				Message: "The current user does not exist",
			},
		},
	}
}

func RequestBodyInvalid() *ApiError {
	return &ApiError{
		StatusCode: http.StatusBadRequest,
		Err: ErrorResponse{
			Error: ErrorDetails{
				Code:    "invalid_request_body",
				Message: "The request body is invalid",
			},
		},
	}
}
