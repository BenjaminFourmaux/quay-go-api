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

func QuayEncryptionKeyNotSet() *ApiError {
	return &ApiError{
		StatusCode: http.StatusInternalServerError,
		Err: ErrorResponse{
			Error: ErrorDetails{
				Code:    "quay_encryption_key_not_set",
				Message: "DATABASE_SECRET_KEY environment variable is not set. Provide the current Quay configuration database secret key.",
			},
		},
	}
}
