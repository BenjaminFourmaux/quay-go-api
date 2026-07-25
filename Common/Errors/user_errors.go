package Errors

import "net/http"

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

func UserNotFound(username string) *ApiError {
	return &ApiError{
		StatusCode: http.StatusNotFound,
		Err: ErrorResponse{
			Error: ErrorDetails{
				Code:    "user_not_found",
				Message: "The user '" + username + "' does not exist",
			},
		},
	}
}
