package Models

type RepositoryPermission struct {
	ID           int  `gorm:"primary_key;auto_increment:false"`
	TeamId       *int `gorm:"type:int;null"`
	UserId       *int `gorm:"type:int;null"`
	RepositoryId int  `gorm:"type:int;not null"`
	RoleId       int  `gorm:"type:int;not null"` // 1 -> admin; 2 -> write; 3 -> read

	// FK
	Team       *Team      `gorm:"foreignkey:TeamId;references:ID"`
	User       *User      `gorm:"foreignkey:UserId;references:ID"`
	Repository Repository `gorm:"foreignkey:RepositoryId;references:ID"`
	Role       Role       `gorm:"foreignkey:RoleId;references:ID"`
}

func (RepositoryPermission) TableName() string {
	return "repositorypermission"
}
