package Errors

type ApiError struct {
	StatusCode int
	Err        ErrorResponse
}

type ErrorResponse struct {
	Error ErrorDetails `json:"error"`
}

type ErrorDetails struct {
	// Code is a sting that represents the error code (e.g. "USER_NOT_FOUND", "INVALID_INPUT", etc.)
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (e *ApiError) Error() string {
	return e.Err.Error.Message
}
