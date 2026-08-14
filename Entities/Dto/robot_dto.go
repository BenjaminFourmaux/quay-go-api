package Dto

import "time"

type Robot struct {
	Name         string     `json:"name"`
	Description  string     `json:"description"`
	Created      time.Time  `json:"created"`
	LastAccessed *time.Time `json:"last_accessed"`
	Token        *string    `json:"token,omitempty"`
	Repositories *[]string  `json:"repositories,omitempty"` // List of repository names the robot has access to
}

type CreateRobot struct {
	Name                 string                 `json:"name"`
	Description          string                 `json:"description"`
	UnstructuredMetadata map[string]interface{} `json:"unstructured_metadata,omitempty"`
}
