package Models

type UserPrompt struct {
	ID     int `gorm:"primaryKey;autoIncrement"`
	UserID int `gorm:"not null"`
	KindId int `gorm:"not null"`

	// FK
	Kind UserPromptKind `gorm:"foreignKey:KindId;references:ID"`
}

func (UserPrompt) TableName() string {
	return "userprompt"
}
