package Models

import "time"

type OauthAccessToken struct {
	ID               int    `gorm:"primaryKey;autoIncrement"`
	UUID             string `gorm:"type:varchar(36);not null"`
	ApplicationID    int    `gorm:"not null"`
	AuthorizedUserID int    `gorm:"not null"`
	//AuthorizedUser   User      `gorm:"foreignKey:AuthorizedUserID"`
	Scope     string    `gorm:"type:varchar(255);not null"` // Scope space-separated list of permissions
	TokenType string    `gorm:"type:varchar(255);not null"`
	ExpiresAt time.Time `gorm:"not null"`
	Data      string    `gorm:"type:text;not null"`
	TokenCode string    `gorm:"type:varchar(255);not null"`
	TokenName string    `gorm:"type:varchar(255);not null"`
}

func (OauthAccessToken) TableName() string {
	return "oauthaccesstoken"
}
