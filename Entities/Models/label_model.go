package Models

type Label struct {
	ID           int    `gorm:"not null;primary_key"`
	UUID         string `gorm:"not null;unique_index"`
	Key          string `gorm:"not null;type:varchar(255)"`
	Value        string `gorm:"not null;type:text"`
	MediaTypeId  int    `gorm:"not null"`
	SourceTypeId int    `gorm:"not null"`

	// FK
	MediaType  MediaType       `gorm:"foreignKey:MediaTypeId;references:ID"`
	SourceType LabelSourceType `gorm:"foreignKey:SourceTypeId;references:ID"`
}

func (Label) TableName() string {
	return "label"
}
