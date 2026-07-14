package Dto

type RepositoryPermission struct {
	Name   string `json:"name"`
	Role   string `json:"role"`
	Avatar Avatar `json:"avatar"`

	// Optional fields
	IsRobot *bool `json:"isRobot,omitempty"`
}

type UpdateRepositoryPermission struct {
	Role string `json:"role"` // allowed value: "admin", "write", "read"
}
