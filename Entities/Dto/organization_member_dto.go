package Dto

type OrganizationMember struct {
	Name         string   `json:"name"`
	Kind         string   `json:"kind"`
	Avatar       Avatar   `json:"avatar"`
	Teams        []string `json:"teams"`        // List of member's team name
	Repositories []string `json:"repositories"` // List of member's repository name
}
