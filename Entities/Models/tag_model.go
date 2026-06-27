package Models

type Tag struct {
	ID              int    `gorm:"primaryKey;autoIncrement"`
	Name            string `gorm:"not null"`
	RepositoryId    int    `gorm:"not null"`
	ManifestId      *int   `gorm:"null"`
	LifetimeStartMs int64  `gorm:"type:bigint;not null"` // Timestamp in milliseconds
	LifetimeEndMs   *int64 `gorm:"type:bigint;null"`     // Timestamp in milliseconds
	Hidden          bool   `gorm:"not null"`
	Revision        int    `gorm:"not null"`
	TagKindId       int    `gorm:"not null"` // 1 -> tag
	LinkedTagId     *int   `gorm:"null"`

	// FK
	Repository Repository `gorm:"foreignKey:RepositoryId;references:ID"`
	Manifest   Manifest   `gorm:"foreignKey:ManifestId;references:ID"`
	TagKind    TagKind    `gorm:"foreignKey:TagKindId;references:ID"`
	LinkedTag  *Tag       `gorm:"foreignKey:LinkedTagId;references:ID"`
}

func (Tag) TableName() string {
	return "tag"
}
