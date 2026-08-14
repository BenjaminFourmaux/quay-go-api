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
