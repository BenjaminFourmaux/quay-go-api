package Errors

import (
	"fmt"
	"net/http"
)

func TagNotFound(tag string, repository string) *ApiError {
	return &ApiError{
		StatusCode: http.StatusBadRequest,
		Err: ErrorResponse{
			Error: ErrorDetails{
				Code:    "tag_not_found",
				Message: fmt.Sprintf("Tag '%s' not found in the repository '%s'", tag, repository),
			},
		},
	}
}
