package Dto

type Team struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Role        string `json:"role"`
	Avatar      Avatar `json:"avatar"`
	CanView     bool   `json:"can_view"`
	//RepoCount    int    `json:"repo_count"` // Useless?
	MembersCount int  `json:"members_count"`
	IsSynced     bool `json:"is_synced"`
}
