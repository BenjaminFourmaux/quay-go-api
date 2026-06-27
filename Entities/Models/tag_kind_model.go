package Models

type TagKind struct {
	ID   uint   `gorm:"primary_key;auto_increment:false"` // 1 -> tag
	Name string `gorm:"not null"`
}

func (TagKind) TableName() string {
	return "tagkind"
}
