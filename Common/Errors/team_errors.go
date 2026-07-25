package Errors

import "net/http"

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
