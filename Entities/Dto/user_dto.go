package Dto

type UserMeResponse struct {
	Anonymous           bool               `json:"anonymous"`
	Username            string             `json:"username"`
	Avatar              Avatar             `json:"avatar"`
	CanCreateRepo       bool               `json:"can_create_repo"`
	IsMe                bool               `json:"is_me"`
	Verified            bool               `json:"verified"`
	Email               string             `json:"email"`
	Logins              []UserLogin        `json:"logins"`
	InvoiceEmail        bool               `json:"invoice_email"`
	InvoiceEmailAddress NullString         `json:"invoice_email_address" swaggertype:"string"`
	PreferredNamespace  bool               `json:"preferred_namespace"`
	TagExpirationS      int                `json:"tag_expiration_s"`
	Prompts             []string           `json:"prompts"` // AI feature in coming ??
	Company             NullString         `json:"company" swaggertype:"string"`
	FamilyName          NullString         `json:"family_name" swaggertype:"string"`
	GivenName           NullString         `json:"given_name" swaggertype:"string"`
	Location            NullString         `json:"location" swaggertype:"string"`
	IsFreeAccount       bool               `json:"is_free_account"`
	HasPasswordSet      bool               `json:"has_password_set"`
	Organizations       []UserOrganization `json:"organizations"`
	SuperUser           bool               `json:"super_user"`
}
