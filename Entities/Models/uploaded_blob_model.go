package Models

import "time"

type UploadedBlob struct {
	ID           int       `gorm:"primaryKey;autoIncrement"`
	RepositoryId int       `gorm:"not null;unique"`
	BlobId       int       `gorm:"not null;unique"`
	UploadedAt   time.Time `gorm:"not null"`
	ExpiresAt    time.Time `gorm:"not null"`

	// FK
	Repository Repository   `gorm:"foreignKey:RepositoryId;references:ID"`
	Blob       ImageStorage `gorm:"foreignKey:BlobId;references:ID"`
}

func (UploadedBlob) TableName() string {
	return "uploadedblob"
}
