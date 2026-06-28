package Errors

import (
	"net/http"
	"quay-go-api/Services/Auth"
	"strings"
)

// <editor-fold desc="Common Errors">

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

// </editor-fold>

// <editor-fold desc="Message Errors">

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

// </editor-fold>

// <editor-fold desc="Organization Errors">

func OrganizationNotFound(orgName string) *ApiError {
	return &ApiError{
		StatusCode: http.StatusNotFound,
		Err: ErrorResponse{
			Error: ErrorDetails{
				Code:    "organization_not_found",
				Message: "The organization '" + orgName + "' does not exist",
			},
		},
	}
}

func UserOrOrganizationNotFound(name string) *ApiError {
	return &ApiError{
		StatusCode: http.StatusNotFound,
		Err: ErrorResponse{
			Error: ErrorDetails{
				Code:    "user_or_organization_not_found",
				Message: "The user or organization '" + name + "' does not exist",
			},
		},
	}
}

func UserOrOrganizationAlreadyExists() *ApiError {
	return &ApiError{
		StatusCode: http.StatusConflict,
		Err: ErrorResponse{
			Error: ErrorDetails{
				Code:    "user_or_organization_already_exists",
				Message: "A user or organization with this name already exists",
			},
		},
	}
}

func OrganizationNameInvalid() *ApiError {
	return &ApiError{
		StatusCode: http.StatusBadRequest,
		Err: ErrorResponse{
			Error: ErrorDetails{
				Code:    "organization_name_invalid",
				Message: "Organization name is invalid. Must be alphanumeric, all lowercase, at least 2 characters long and at most 255 characters long",
			},
		},
	}
}

func OrganizationEmailInvalid() *ApiError {
	return &ApiError{
		StatusCode: http.StatusBadRequest,
		Err: ErrorResponse{
			Error: ErrorDetails{
				Code:    "organization_email_invalid",
				Message: "Organization email is invalid",
			},
		},
	}
}

func OrganizationTagExpirationInvalid() *ApiError {
	return &ApiError{
		StatusCode: http.StatusBadRequest,
		Err: ErrorResponse{
			Error: ErrorDetails{
				Code:    "organization_tag_expiration_invalid",
				Message: "Tag expiration must be greater than or equal to 0",
			},
		},
	}
}

func UserNotOrganizationOwner() *ApiError {
	return &ApiError{
		StatusCode: http.StatusForbidden,
		Err: ErrorResponse{
			Error: ErrorDetails{
				Code:    "user_not_organization_owner",
				Message: "The user is not the owner of this organization",
			},
		},
	}
}

// </editor-fold>

// <editor-fold desc="Team Errors">

func TeamNotFound(teamName string) *ApiError {
	return &ApiError{
		StatusCode: http.StatusNotFound,
		Err: ErrorResponse{
			Error: ErrorDetails{
				Code:    "team_not_found",
				Message: "The team '" + teamName + "' does not exist",
			},
		},
	}
}

func TeamAlreadyExists() *ApiError {
	return &ApiError{
		StatusCode: http.StatusConflict,
		Err: ErrorResponse{
			Error: ErrorDetails{
				Code:    "team_already_exists",
				Message: "A team with this name already exists",
			},
		},
	}
}

func TeamNameRequired() *ApiError {
	return &ApiError{
		StatusCode: http.StatusBadRequest,
		Err: ErrorResponse{
			Error: ErrorDetails{
				Code:    "team_name_required",
				Message: "Team name is required",
			},
		},
	}
}

func TeamNameInvalid() *ApiError {
	return &ApiError{
		StatusCode: http.StatusBadRequest,
		Err: ErrorResponse{
			Error: ErrorDetails{
				Code:    "team_name_invalid",
				Message: "Team name is invalid. Must be alphanumeric, all lowercase, at least 2 characters long and at most 255 characters long",
			},
		},
	}
}

func TeamRoleInvalid() *ApiError {
	return &ApiError{
		StatusCode: http.StatusBadRequest,
		Err: ErrorResponse{
			Error: ErrorDetails{
				Code:    "team_role_invalid",
				Message: "Team role is invalid. Must be one of 'member', 'admin' or 'creator'",
			},
		},
	}
}

// </editor-fold>

// <editor-fold desc="Members Errors">

func MemberAlreadyInTeam() *ApiError {
	return &ApiError{
		StatusCode: http.StatusConflict,
		Err: ErrorResponse{
			Error: ErrorDetails{
				Code:    "member_already_in_team",
				Message: "The user is already a member of this team",
			},
		},
	}
}

func MemberNotInTeam() *ApiError {
	return &ApiError{
		StatusCode: http.StatusConflict,
		Err: ErrorResponse{
			Error: ErrorDetails{
				Code:    "member_not_in_team",
				Message: "The user is not a member of this team",
			},
		},
	}
}

// </editor-fold>

// <editor-fold desc="User Errors">

func UserNotExists(username string) *ApiError {
	return &ApiError{
		StatusCode: http.StatusNotFound,
		Err: ErrorResponse{
			Error: ErrorDetails{
				Code:    "user_not_exists",
				Message: "The user '" + username + "' does not exist",
			},
		},
	}
}

// </editor-fold>

// <editor-fold desc="Repository Errors">

func RepositoryNameInvalid() *ApiError {
	return &ApiError{
		StatusCode: http.StatusBadRequest,
		Err: ErrorResponse{
			Error: ErrorDetails{
				Code:    "repository_name_invalid",
				Message: "Repository name is invalid. It must use lowercase letters and digits, may include separators ('.', '_', '__', or '-'), may use '/' between name components, each component must start and end with a letter or digit, and the full name must be at most 255 characters long",
			},
		},
	}
}

func RepositoryKindInvalid() *ApiError {
	return &ApiError{
		StatusCode: http.StatusBadRequest,
		Err: ErrorResponse{
			Error: ErrorDetails{
				Code:    "repository_invalid_kind",
				Message: "Repository kind is invalid. Valid kinds are 'image' and 'application'",
			},
		},
	}
}

func RepositoryNamespaceInvalid() *ApiError {
	return &ApiError{
		StatusCode: http.StatusBadRequest,
		Err: ErrorResponse{
			Error: ErrorDetails{
				Code:    "repository_namespace_invalid",
				Message: "Repository namespace is invalid. Must be alphanumeric, all lowercase, at least 2 characters long and at most 255 characters long",
			},
		},
	}
}

func RepositoryInvalid(repository string) *ApiError {
	return &ApiError{
		StatusCode: http.StatusBadRequest,
		Err: ErrorResponse{
			Error: ErrorDetails{
				Code:    "repository_invalid",
				Message: "Invalid repository '" + repository + "'",
			},
		},
	}
}

func RepositoryNamespaceNotFound(namespace string) *ApiError {
	return &ApiError{
		StatusCode: http.StatusNotFound,
		Err: ErrorResponse{
			Error: ErrorDetails{
				Code:    "repository_namespace_not_found",
				Message: "The namespace '" + namespace + "' does not exist. Provide a valid user or organization name",
			},
		},
	}
}

func RepositoryNotFound(repository string) *ApiError {
	return &ApiError{
		StatusCode: http.StatusNotFound,
		Err: ErrorResponse{
			Error: ErrorDetails{
				Code:    "repository_not_found",
				Message: "The repository '" + repository + "' does not exist",
			},
		},
	}
}

func RepositoryAlreadyExists() *ApiError {
	return &ApiError{
		StatusCode: http.StatusBadRequest,
		Err: ErrorResponse{
			Error: ErrorDetails{
				Code:    "repository_already_exists",
				Message: "Repository already exists",
			},
		},
	}
}

// </editor-fold>
