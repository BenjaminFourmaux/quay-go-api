package Errors

import (
	"fmt"
	"net/http"
)

func PermissionNotFound(kind string, name string) *ApiError {
	return &ApiError{
		StatusCode: http.StatusNotFound,
		Err: ErrorResponse{
			Error: ErrorDetails{
				Code:    "permission_not_found",
				Message: fmt.Sprintf("Permission not found for '%s': '%s'", kind, name),
			},
		},
	}
}

func RepositoryPermissionRoleInvalid(role string) *ApiError {
	return &ApiError{
		StatusCode: http.StatusBadRequest,
		Err: ErrorResponse{
			Error: ErrorDetails{
				Code:    "repository_permission_role_invalid",
				Message: fmt.Sprintf("Repository permission role '%s' is invalid. Valid roles are 'admin', 'write', and 'read'", role),
			},
		},
	}
}
