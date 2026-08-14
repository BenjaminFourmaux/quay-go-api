package Errors

import "net/http"

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
