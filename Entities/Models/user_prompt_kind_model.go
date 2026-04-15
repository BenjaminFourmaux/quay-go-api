package Models

type UserPromptKind struct {
	ID   int    `gorm:"primaryKey;autoIncrement"`
	Name string `gorm:"type:varchar(255);not null"`
}

func (UserPromptKind) TableName() string {
	return "userpromptkind"
}
