package Models

type TeamRole struct {
	ID   string `gorm:"type:int;primaryKey;autoIncrement"`
	Name string `gorm:"type:varchar(255);not null"`
}

func (TeamRole) TableName() string {
	return "teamrole"
}
