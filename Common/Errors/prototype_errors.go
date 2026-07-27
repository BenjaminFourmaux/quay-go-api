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
