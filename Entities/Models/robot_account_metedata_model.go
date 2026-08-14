package Models

type RobotAccountMetadata struct {
	ID               int    `gorm:"primaryKey;autoIncrement"`
	RobotAccountID   int    `gorm:"not null"`
	Description      string `gorm:"not null"`
	UnstructuredJson string `gorm:"type:text"` // Store unstructured data as JSONB

	// FK
	RobotAccount User `gorm:"foreignKey:RobotAccountID;references:ID"`
}

func (r *RobotAccountMetadata) TableName() string {
	return "robotaccountmetadata"
}
