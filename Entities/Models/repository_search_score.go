package Models

import "time"

type RepositorySearchScore struct {
	ID           int        `gorm:"primary_key"`
	RepositoryId int        `gorm:"index"`
	Score        int64      `gorm:"type:bigint;index"`
	LastUpdated  *time.Time `gorm:"type:timestamp without time zone;default:null"`

	// FK
	Repository *Repository `gorm:"foreignKey:RepositoryId;references:Id"`
}

func (f *RepositorySearchScore) TableName() string {
	return "repositorysearchscore"
}
