package Models

type PermissionPrototype struct {
	ID               int    `gorm:"primary_key"`
	UUID             string `gorm:"not null;unique"`
	OrgId            int    `gorm:"not null"`
	ActivatingUserId *int   `gorm:"null"`
	DelegateUserId   *int   `gorm:"null"`
	DelegateTeamId   *int   `gorm:"null"`
	RoleId           int    `gorm:"not null"`

	// FK
	Organization   User  `gorm:"foreignKey:OrgId;references:ID"`
	ActivatingUser *User `gorm:"foreignKey:ActivatingUserId;references:ID"`
	DelegateUser   *User `gorm:"foreignKey:DelegateUserId;references:ID"`
	DelegateTeam   *Team `gorm:"foreignKey:DelegateTeamId;references:ID"`
	Role           Role  `gorm:"foreignKey:RoleId;references:ID"`
}

func (*PermissionPrototype) TableName() string {
	return "permissionprototype"
}
