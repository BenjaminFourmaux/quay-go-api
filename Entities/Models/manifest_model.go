package Models

type Manifest struct {
	ID                   uint    `gorm:"primary_key;not null;column:id"`
	RepositoryId         int     `gorm:"primary_key;not null;column:repository_id"`
	Digest               string  `gorm:"primary_key;not null;column:digest"`
	MediaTypeId          int     `gorm:"type:int;not null;column:media_type_id"`
	ManifestBytes        string  `gorm:"type:text;not null;column:manifest_bytes"`
	ConfigMediaType      *string `gorm:"type:varchar(255);null;column:config_media_type"`
	LayersCompressedSize *int64  `gorm:"type:bigint;null;column:layers_compressed_size"` // In bits

	// FK
	Repository Repository `gorm:"foreignKey:RepositoryId;references:ID"`
	MediaType  MediaType  `gorm:"foreignKey:MediaTypeId;references:ID"`
}

func (Manifest) TableName() string {
	return "manifest"
}
