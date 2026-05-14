package Models

type Team struct {
	ID             int    `gorm:"primaryKey;autoIncrement"`
	Name           string `gorm:"type:varchar(255);not null"`
	Description    string `gorm:"type:varchar(255);not null"`
	OrganizationId int    `gorm:"type:int;not null"`
	RoleId         int    `gorm:"type:int;not null"`

	// FK
	Organization User         `gorm:"foreignKey:OrganizationId;references:ID"`
	Role         TeamRole     `gorm:"foreignKey:RoleId;references:ID"`
	Members      []TeamMember `gorm:"foreignKey:TeamId;references:ID"`
}

func (Team) TableName() string {
	return "team"
}
