package Models

import "time"

type Star struct {
	ID           int       `gorm:"type:int;primaryKey;autoIncrement"`
	UserId       int       `gorm:"type:varchar(255);not null"`
	RepositoryId int       `gorm:"type:int;not null"`
	Created      time.Time `gorm:"type:timestamp;not null"`

	// FK
	User       User       `gorm:"foreignKey:UserId;references:Id"`
	Repository Repository `gorm:"foreignKey:RepositoryId;references:Id"`
}
