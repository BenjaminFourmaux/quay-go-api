package Models

type TeamMember struct {
	ID     int `gorm:"primary_key"`
	UserId int
	TeamId int

	// FK
	User User `gorm:"foreignKey:UserId;references:ID"`
	Team Team `gorm:"foreignKey:TeamId;references:ID"`
}

func (TeamMember) TableName() string {
	return "teammember"
}
