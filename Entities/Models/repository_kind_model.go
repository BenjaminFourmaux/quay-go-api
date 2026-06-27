package Models

type RepositoryKind struct {
	ID   int    `gorm:"primary_key;auto_increment:false"` // 1 = image , 2 = application
	Name string `gorm:"type:varchar(255);not null"`
}

func (RepositoryKind) TableName() string {
	return "repositorykind"
}
