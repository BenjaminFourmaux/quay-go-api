package Dto

type Prototype struct {
	Id             string         `json:"id"` // UUID of the prototype
	Role           string         `json:"role"`
	ActivatingUser ActivatingUser `json:"activating_user"`
	Delegate       Delegate       `json:"delegate"`
}

type ActivatingUser struct {
	Name        string `json:"name"`
	IsRobot     bool   `json:"is_robot"`
	Kind        string `json:"kind"` // 'user' or 'team'
	IsOrgMember bool   `json:"is_org_member"`
	Avatar      Avatar `json:"avatar"`
}

type Delegate struct {
	Name        string `json:"name"`
	IsRobot     bool   `json:"is_robot"`
	Kind        string `json:"kind"` // 'user' or 'team'
	IsOrgMember bool   `json:"is_org_member"`
	Avatar      Avatar `json:"avatar"`
}
