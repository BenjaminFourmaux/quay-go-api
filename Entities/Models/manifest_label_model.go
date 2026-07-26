package Models

/*
ManifestLabel Jointure table
*/
type ManifestLabel struct {
	ID           int `gorm:"primary_key"`
	RepositoryId int `gorm:"index"`
	ManifestId   int `gorm:"index"`
	LabelId      int `gorm:"index"`

	// FK
	Repository Repository `gorm:"foreignKey:RepositoryId;references:ID"`
	Manifest   Manifest   `gorm:"foreignKey:ManifestId;references:ID"`
	Label      Label      `gorm:"foreignKey:LabelId;references:ID"`
}

func (ManifestLabel) TableName() string {
	return "manifestlabel"
}
