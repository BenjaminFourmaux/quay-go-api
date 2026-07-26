package Errors

import (
	"fmt"
	"net/http"
)

func QuayUrlNotSet() *ApiError {
	return &ApiError{
		StatusCode: http.StatusInternalServerError,
		Err: ErrorResponse{
			Error: ErrorDetails{
				Code:    "quay_url_not_set",
				Message: "QUAY_URL environment variable is not set",
			},
		},
	}
}

func QuayApiError(errorCode int, errMessage string) *ApiError {
	return &ApiError{
		StatusCode: http.StatusInternalServerError,
		Err: ErrorResponse{
			Error: ErrorDetails{
				Code:    "quay_api_error",
				Message: fmt.Sprintf("[HTTP %d] %s", errorCode, errMessage),
			},
		},
	}
}
