package Errors

import (
	"net/http"
	"strconv"
)

func RobotNameRequired() *ApiError {
	return &ApiError{
		StatusCode: http.StatusBadRequest,
		Err: ErrorResponse{
			Error: ErrorDetails{
				Code:    "robot_name_required",
				Message: "Robot name is required",
			},
		},
	}
}

func RobotNameInvalid(regexp string) *ApiError {
	return &ApiError{
		StatusCode: http.StatusBadRequest,
		Err: ErrorResponse{
			Error: ErrorDetails{
				Code:    "robot_name_invalid",
				Message: "Robot name is invalid. It must follow the pattern: " + regexp,
			},
		},
	}
}

func RobotAlreadyExists(name string) *ApiError {
	return &ApiError{
		StatusCode: http.StatusConflict,
		Err: ErrorResponse{
			Error: ErrorDetails{
				Code:    "robot_already_exists",
				Message: "A robot with the name '" + name + "' already exists",
			},
		},
	}
}

func RobotNotFound(name string) *ApiError {
	return &ApiError{
		StatusCode: http.StatusNotFound,
		Err: ErrorResponse{
			Error: ErrorDetails{
				Code:    "robot_not_found",
				Message: "No robot found with the name '" + name + "'",
			},
		},
	}
}

func RobotNotFoundById(robotId int) *ApiError {
	return &ApiError{
		StatusCode: http.StatusNotFound,
		Err: ErrorResponse{
			Error: ErrorDetails{
				Code:    "robot_not_found",
				Message: "No robot found with the id '" + strconv.Itoa(robotId) + "'",
			},
		},
	}
}

func FederationInvalidIssuer() *ApiError {
	return &ApiError{
		StatusCode: http.StatusBadRequest,
		Err: ErrorResponse{
			Error: ErrorDetails{
				Code:    "federation_invalid_issuer",
				Message: "Required Federation issuer is invalid or missing",
			},
		},
	}
}

func FederationInvalidSubject() *ApiError {
	return &ApiError{
		StatusCode: http.StatusBadRequest,
		Err: ErrorResponse{
			Error: ErrorDetails{
				Code:    "federation_invalid_subject",
				Message: "Required Federation subject is invalid or missing",
			},
		},
	}
}

func FederationIssuerMustBeURL() *ApiError {
	return &ApiError{
		StatusCode: http.StatusBadRequest,
		Err: ErrorResponse{
			Error: ErrorDetails{
				Code:    "federation_issuer_must_be_url",
				Message: "Issuer must be a URL (http:// or https://)",
			},
		},
	}
}
