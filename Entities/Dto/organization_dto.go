package Dto

// TODO: maybe rationalized this interface for make common (ex: /user/me get the same orgs detail than /organization/{orgname} ??

type Organization struct {
	Name                string     `json:"name"`
	Avatar              Avatar     `json:"avatar"`
	IsAdmin             bool       `json:"is_admin"`
	IsMember            bool       `json:"is_member"`
	Teams               []Team     `json:"teams"`
	InvoiceEmail        bool       `json:"invoice_email"`
	InvoiceEmailAddress NullString `json:"invoice_email_address" swaggertype:"string"`
	TagExpirationS      int        `json:"tag_expiration_s"`
	IsFreeAccount       bool       `json:"is_free_account"`
}

type UserOrganization struct {
	Name               string `json:"name"`
	Avatar             Avatar `json:"avatar"`
	CanCreateRepo      bool   `json:"can_create_repo"`
	Public             bool   `json:"public"`
	IsOrgAdmin         bool   `json:"is_org_admin"`
	PreferredNamespace bool   `json:"preferred_namespace"`
}
