package Errors

import (
	"fmt"
	"net/http"
)

func PrototypeDelegateKindInvalid(kind string) *ApiError {
	return &ApiError{
		StatusCode: http.StatusBadRequest,
		Err: ErrorResponse{
			Error: ErrorDetails{
				Code:    "prototype_delegate_kind_invalid",
				Message: fmt.Sprintf("Prototype delegate kind '%s' is invalid. Valid kinds are 'user' and 'team'", kind),
			},
		},
	}
}

func PrototypeAlreadyExists(uuid string) *ApiError {
	return &ApiError{
		StatusCode: http.StatusConflict,
		Err: ErrorResponse{
			Error: ErrorDetails{
				Code:    "prototype_already_exists",
				Message: fmt.Sprintf("The prototype already exists with id '%s'", uuid),
			},
		},
	}
}

func PrototypeNotFound(uuid string) *ApiError {
	return &ApiError{
		StatusCode: http.StatusNotFound,
		Err: ErrorResponse{
			Error: ErrorDetails{
				Code:    "prototype_not_found",
				Message: fmt.Sprintf("The prototype '%s' not found", uuid),
			},
		},
	}
}
