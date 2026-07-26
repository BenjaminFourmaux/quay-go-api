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

func TagNameInvalid(tag string) *ApiError {
	return &ApiError{
		StatusCode: http.StatusBadRequest,
		Err: ErrorResponse{
			Error: ErrorDetails{
				Code:    "tag_name_invalid",
				Message: fmt.Sprintf("Tag '%s' is invalid. It must match the pattern '%s'", tag, "^[\\w][\\w.-]{0,127}$"),
			},
		},
	}
}

func ManifestNotFound(manifestDigest string, repository string) *ApiError {
	return &ApiError{
		StatusCode: http.StatusNotFound,
		Err: ErrorResponse{
			Error: ErrorDetails{
				Code:    "manifest_not_found",
				Message: fmt.Sprintf("Manifest '%s' not found in the repository '%s'", manifestDigest, repository),
			},
		},
	}
}

func InvalidExpirationDate(expirationDate string) *ApiError {
	return &ApiError{
		StatusCode: http.StatusBadRequest,
		Err: ErrorResponse{
			Error: ErrorDetails{
				Code:    "tag_expiration_invalid",
				Message: fmt.Sprintf("Expiration date '%s' is invalid. It must be a future date.", expirationDate),
			},
		},
	}
}

func ManifestLabelNotFound(labelId string, manifestDigest string, repository string) *ApiError {
	return &ApiError{
		StatusCode: http.StatusNotFound,
		Err: ErrorResponse{
			Error: ErrorDetails{
				Code:    "manifest_label_not_found",
				Message: fmt.Sprintf("Manifest label '%s' not found for manifest '%s' in the repository '%s'", labelId, manifestDigest, repository),
			},
		},
	}
}
