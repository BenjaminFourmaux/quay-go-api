package Models

type ImageStorage struct {
	ID               int    `gorm:"primaryKey;autoIncrement"`
	UUID             string `gorm:"not null;unique"`
	ImageSize        int    `gorm:"null"`
	UncompressedSize int    `gorm:"null"`
	Uploading        bool   `gorm:"null"`
	CasPath          bool   `gorm:"not null;default:false"`
	ContentChecksum  string `gorm:"null"`
}

func (ImageStorage) TableName() string {
	return "imagestorage"
}
