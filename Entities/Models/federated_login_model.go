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

/*
FederatedLoginMetadata Model for serialize/deserialize json from field FederationLogin.MetadataJson
*/
type FederatedLoginMetadata struct {
	FederationConfig []FederationConfig `json:"federation_config"`
}
type FederationConfig struct {
	Issuer  string `json:"issuer"`
	Subject string `json:"subject"`
}
