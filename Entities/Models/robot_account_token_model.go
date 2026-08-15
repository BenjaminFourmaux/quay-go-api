package Models

type RobotAccountToken struct {
	ID             int    `gorm:"primaryKey;autoIncrement"`
	RobotAccountID int    `gorm:"not null"`
	Token          string `gorm:"not null"` // Encrypted token value
	FullyMigrated  bool   `gorm:"type:boolean"`
	// FK
	RobotAccount User `gorm:"foreignKey:RobotAccountID;references:ID"`
}

func (r *RobotAccountToken) TableName() string {
	return "robotaccounttoken"
}
