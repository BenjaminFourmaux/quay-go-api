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

type CreateTeam struct {
	Name        *string `json:"name"`
	Description *string `json:"description"` // Description of the Team can be in Markdown format (optional)
	Role        *string `json:"role"`        // Name of the role ('admin', 'creator' or 'member'. Optional. Default 'member')
}

type UpdateTeam struct {
	Description *string `json:"description"`
	Role        *string `json:"role"`
}
