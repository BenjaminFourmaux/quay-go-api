package Models

import "time"

type RepositoryActionCount struct {
	ID           int       `gorm:"type:int;primaryKey;autoIncrement"`
	RepositoryId int       `gorm:"type:int;not null"`
	Count        int       `gorm:"type:int;not null"`
	Date         time.Time `gorm:"type:date;not null"`

	// FK
	Repository Repository `gorm:"foreignKey:RepositoryId;references:Id"`
}

func (f *RepositoryActionCount) TableName() string {
	return "repositoryactioncount"
}
