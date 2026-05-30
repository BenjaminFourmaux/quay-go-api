package Models

type TeamMemberInvite struct {
	ID          int    `gorm:"primaryKey;autoIncrement"`
	UserId      int    `gorm:"type:int;null"`
	Email       string `gorm:"type:string;null"`
	TeamId      int    `gorm:"type:int;not null"`
	InviterId   int    `gorm:"type:int;not null"`
	InviteToken string `gorm:"type:string;not null"`

	// FK
	User    User `gorm:"foreignKey:UserId;references:ID"`
	Team    Team `gorm:"foreignKey:TeamId;references:ID"`
	Inviter User `gorm:"foreignKey:InviterId;references:ID"`
}

func (TeamMemberInvite) TableName() string {
	return "teammemberinvite"
}
