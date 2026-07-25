package Errors

import "net/http"

func OrganizationNotFound(orgName string) *ApiError {
	return &ApiError{
		StatusCode: http.StatusNotFound,
		Err: ErrorResponse{
			Error: ErrorDetails{
				Code:    "organization_not_found",
				Message: "The organization '" + orgName + "' does not exist",
			},
		},
	}
}

func UserOrOrganizationNotFound(name string) *ApiError {
	return &ApiError{
		StatusCode: http.StatusNotFound,
		Err: ErrorResponse{
			Error: ErrorDetails{
				Code:    "user_or_organization_not_found",
				Message: "The user or organization '" + name + "' does not exist",
			},
		},
	}
}

func UserOrOrganizationAlreadyExists() *ApiError {
	return &ApiError{
		StatusCode: http.StatusConflict,
		Err: ErrorResponse{
			Error: ErrorDetails{
				Code:    "user_or_organization_already_exists",
				Message: "A user or organization with this name already exists",
			},
		},
	}
}

func OrganizationNameInvalid() *ApiError {
	return &ApiError{
		StatusCode: http.StatusBadRequest,
		Err: ErrorResponse{
			Error: ErrorDetails{
				Code:    "organization_name_invalid",
				Message: "Organization name is invalid. Must be alphanumeric, all lowercase, at least 2 characters long and at most 255 characters long",
			},
		},
	}
}

func OrganizationEmailInvalid() *ApiError {
	return &ApiError{
		StatusCode: http.StatusBadRequest,
		Err: ErrorResponse{
			Error: ErrorDetails{
				Code:    "organization_email_invalid",
				Message: "Organization email is invalid",
			},
		},
	}
}

func OrganizationTagExpirationInvalid() *ApiError {
	return &ApiError{
		StatusCode: http.StatusBadRequest,
		Err: ErrorResponse{
			Error: ErrorDetails{
				Code:    "organization_tag_expiration_invalid",
				Message: "Tag expiration must be greater than or equal to 0",
			},
		},
	}
}

func UserNotOrganizationOwner() *ApiError {
	return &ApiError{
		StatusCode: http.StatusForbidden,
		Err: ErrorResponse{
			Error: ErrorDetails{
				Code:    "user_not_organization_owner",
				Message: "The user is not the owner of this organization",
			},
		},
	}
}
