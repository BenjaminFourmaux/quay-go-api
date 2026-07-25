package Errors

import "net/http"

func RepositoryNameInvalid() *ApiError {
	return &ApiError{
		StatusCode: http.StatusBadRequest,
		Err: ErrorResponse{
			Error: ErrorDetails{
				Code:    "repository_name_invalid",
				Message: "Repository name is invalid. It must use lowercase letters and digits, may include separators ('.', '_', '__', or '-'), may use '/' between name components, each component must start and end with a letter or digit, and the full name must be at most 255 characters long",
			},
		},
	}
}

func RepositoryKindInvalid() *ApiError {
	return &ApiError{
		StatusCode: http.StatusBadRequest,
		Err: ErrorResponse{
			Error: ErrorDetails{
				Code:    "repository_invalid_kind",
				Message: "Repository kind is invalid. Valid kinds are 'image' and 'application'",
			},
		},
	}
}

func RepositoryNamespaceInvalid() *ApiError {
	return &ApiError{
		StatusCode: http.StatusBadRequest,
		Err: ErrorResponse{
			Error: ErrorDetails{
				Code:    "repository_namespace_invalid",
				Message: "Repository namespace is invalid. Must be alphanumeric, all lowercase, at least 2 characters long and at most 255 characters long",
			},
		},
	}
}

func RepositoryInvalid(repository string) *ApiError {
	return &ApiError{
		StatusCode: http.StatusBadRequest,
		Err: ErrorResponse{
			Error: ErrorDetails{
				Code:    "repository_invalid",
				Message: "Invalid repository '" + repository + "'",
			},
		},
	}
}

func RepositoryNamespaceNotFound(namespace string) *ApiError {
	return &ApiError{
		StatusCode: http.StatusNotFound,
		Err: ErrorResponse{
			Error: ErrorDetails{
				Code:    "repository_namespace_not_found",
				Message: "The namespace '" + namespace + "' does not exist. Provide a valid user or organization name",
			},
		},
	}
}

func RepositoryNotFound(repository string) *ApiError {
	return &ApiError{
		StatusCode: http.StatusNotFound,
		Err: ErrorResponse{
			Error: ErrorDetails{
				Code:    "repository_not_found",
				Message: "The repository '" + repository + "' does not exist",
			},
		},
	}
}

func RepositoryAlreadyExists() *ApiError {
	return &ApiError{
		StatusCode: http.StatusBadRequest,
		Err: ErrorResponse{
			Error: ErrorDetails{
				Code:    "repository_already_exists",
				Message: "Repository already exists",
			},
		},
	}
}
