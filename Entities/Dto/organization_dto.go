package Dto

type Organization struct {
	Name               string `json:"name"`
	Avatar             Avatar `json:"avatar"`
	CanCreateRepo      bool   `json:"can_create_repo"`
	Public             bool   `json:"public"`
	IsOrgAdmin         bool   `json:"is_org_admin"`
	PreferredNamespace bool   `json:"preferred_namespace"`
}
