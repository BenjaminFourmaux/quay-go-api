package Models

type FederatedLogin struct {
	ID           int    `gorm:"primaryKey;autoIncrement"`
	UserId       int    `gorm:"not null"`
	ServiceId    int    `gorm:"not null"`
	ServiceIdent string `gorm:"not null"`
	MetadataJson string `gorm:"type:text"`

	// FK
	Service LoginService `gorm:"foreignKey:ServiceId;references:ID"`
}

func (f *FederatedLogin) TableName() string {
	return "federatedlogin"
}
