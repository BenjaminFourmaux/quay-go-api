package Repositories

import (
	"gorm.io/gorm"
	"quay-go-api/Database"
	"quay-go-api/Entities/Models"
)

func GetRepositoryManifestByDigest(repositoryId int, digest string) (Models.Manifest, error) {
	var manifest Models.Manifest
	err := Database.DB.
		Where("repository_id = ?", repositoryId).
		Where("digest = ?", digest).
		First(&manifest).Error
	return manifest, err
}

func AddManifestLabel(repoId int, manifestId int, label Models.Label) (Models.Label, error) {
	err := Database.DB.Transaction(func(tx *gorm.DB) error {
		// 1. Insert the Label
		if err := tx.Create(&label).Error; err != nil {
			return err
		}

		// 2. Create jointure object
		manifestLabelAssociation := Models.ManifestLabel{
			RepositoryId: repoId,
			ManifestId:   manifestId,
			LabelId:      label.ID,
		}

		// 2. Insert the ManifestLabel association
		if err := tx.Create(&manifestLabelAssociation).Error; err != nil {
			return err
		}
		return nil
	})

	return label, err
}
