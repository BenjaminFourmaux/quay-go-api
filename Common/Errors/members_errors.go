package Errors

import "net/http"

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
