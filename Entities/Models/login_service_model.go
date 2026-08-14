package Models

type LoginService struct {
	ID   int    `gorm:"primaryKey;autoIncrement"` // 1 => github, 2 => quayrobot, 3 => ldap, 4 => google, 5 => keystone, 6 => dex, 7 => jwtauthn
	Name string `gorm:"not null;unique"`
}

func (l *LoginService) TableName() string {
	return "loginservice"
}
