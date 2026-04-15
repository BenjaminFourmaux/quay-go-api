package Models

type MediaType struct {
	Id   int    `gorm:"primaryKey;autoIncrement"`
	Name string `gorm:"not null;unique"` // Name of the media type (text, markdown, etc.)
}

func (m *MediaType) TableName() string {
	return "mediatype"
}
