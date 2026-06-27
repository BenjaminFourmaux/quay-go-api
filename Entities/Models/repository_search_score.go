package Models

import "time"

type RepositorySearchScore struct {
	ID           int        `gorm:"primary_key;not null;column:id"`
	RepositoryId int        `gorm:"type:int;not null;column:repository_id"`
	Score        int64      `gorm:"type:bigint;not null;column:score"`
	LastUpdated  *time.Time `gorm:"type:timestamp without time zone;default:null;column:last_updated"`

	// FK
	Repository *Repository `gorm:"foreignKey:RepositoryId;references:Id"`
}

func (f *RepositorySearchScore) TableName() string {
	return "repositorysearchscore"
}
