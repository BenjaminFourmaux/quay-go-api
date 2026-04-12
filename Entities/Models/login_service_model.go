package Models

type LoginService struct {
	ID   int    `gorm:"primaryKey;autoIncrement"`
	Name string `gorm:"not null;unique"`
}

func (l *LoginService) TableName() string {
	return "loginservice"
}
