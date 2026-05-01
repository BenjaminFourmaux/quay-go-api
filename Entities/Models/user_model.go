package Models

import (
	"database/sql"
	"time"
)

type User struct {
	ID                       int            `gorm:"primaryKey;autoIncrement"`
	UUID                     string         `gorm:"type:varchar(255);not null"`
	Username                 string         `gorm:"type:varchar(255);not null"`
	PasswordHash             sql.NullString `gorm:"type:varchar(255);null"`               // Can be null for Org user
	Email                    string         `gorm:"type:varchar(255);not null"`           // Can be a UUID for Org user
	Verified                 bool           `gorm:"not null"`                             // if the user is verified or not
	StripeId                 sql.NullString `gorm:"type:varchar(255);null"`               // Stripe customer ID, can be null
	Organization             bool           `gorm:"not null"`                             // Determines if the user is an organization or a regular user (idk why Orgs are users !?)
	Robot                    bool           `gorm:"not null"`                             // Determines if the user is a robot or a regular user (okay, that makes sense for robot account)
	InvoiceEmail             bool           `gorm:"not null"`                             // if the user has an invoice email
	LastInvalidLogin         time.Time      `gorm:"type:time without time zone;not null"` // Last invalid login time
	RemovedTagExpirationS    int            `gorm:"not null"`                             // Expiration time in seconds for removed tags, default is 1209600 (14 days)
	Enabled                  bool           `gorm:"not null"`                             // If the user is enabled or not
	InvoiceEmailAddress      sql.NullString `gorm:"type:varchar(255);null"`               // Invoice email address, can be null
	Company                  sql.NullString `gorm:"type:varchar(255);null"`               // Company name, can be null
	FamilyName               sql.NullString `gorm:"type:varchar(255);null"`               // Family name, can be null
	GivenName                sql.NullString `gorm:"type:varchar(255);null"`               // Given name, can be null
	Location                 sql.NullString `gorm:"type:varchar(255);null"`               // Location, can be null
	MaximumQueuedBuildsCount sql.NullInt16  `gorm:"null"`                                 // Maximum number of queued builds, can be null
	CreationDate             sql.NullTime   `gorm:"type:time without time zone;null"`     // Creation date, can be null
	LastAccessed             sql.NullTime   `gorm:"type:time without time zone;null"`     // Last accessed time, can be null

	// Fk
	FederatedLogins []FederatedLogin `gorm:"foreignKey:UserId;references:ID"`
	Prompts         []UserPrompt     `gorm:"foreignKey:UserId;references:ID"`
	Teams           []Team           `gorm:"foreignKey:OrganizationId;references:ID"` // Only for Organization
}

func (User) TableName() string {
	return "user"
}
