package Models

type Repository struct {
	ID              int    `gorm:"primaryKey;autoIncrement"`
	NamespaceUserId int    `gorm:"type:int;null"`
	Name            string `gorm:"type:varchar(255);not null"`
	VisibilityId    int    `gorm:"type:int;not null"`
	Description     string `gorm:"type:varchar(255);null"`
	BadgeToken      string `gorm:"type:varchar(255);not null"`
	KindId          int    `gorm:"type:int;not null"`
	TrustEnabled    bool   `gorm:"type:bool;not null"`
	State           int    `gorm:"column:state;type:int;not null"` // Enum: NORMAL = 0, READ_ONLY = 1, MIRROR = 2 MARKED_FOR_DELETION = 3, ORG_MIRROR = 4

	// FK
	NamespaceUser User           `gorm:"foreignKey:NamespaceUserId;references:ID"`
	Kind          RepositoryKind `gorm:"foreignKey:KindId;references:ID"`
	Visibility    Visibility     `gorm:"foreignKey:VisibilityId;references:ID"`
	Stars         []Star         `gorm:"foreignKey:RepositoryId;references:ID"`
}

func (f *Repository) TableName() string {
	return "repository"
}
