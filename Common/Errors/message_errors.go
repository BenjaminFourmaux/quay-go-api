package Errors

import "net/http"

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
