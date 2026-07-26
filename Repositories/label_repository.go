package Repositories

import (
	"gorm.io/gorm"
	"quay-go-api/Database"
	"quay-go-api/Entities/Models"
)

func ListManifestLabels(repositoryId int, manifestId int) ([]Models.ManifestLabel, error) {
	var labels []Models.ManifestLabel
	err := Database.DB.
		Preload("Label").
		Where("repository_id = ? AND manifest_id = ?", repositoryId, manifestId).
		Find(&labels).
		Error
	return labels, err
}

func GetManifestLabelByUUID(repositoryId int, manifestId int, labelUUID string) (*Models.ManifestLabel, error) {
	var label Models.ManifestLabel
	err := Database.DB.
		Preload("Label").
		InnerJoins("JOIN label ON label.id = manifestlabel.label_id").
		Where("repository_id = ? AND manifest_id = ? AND label.uuid = ?", repositoryId, manifestId, labelUUID).
		First(&label).
		Error
	return &label, err
}

func DeleteManifestLabelByUUID(repositoryId int, manifestId int, labelUUID string) error {
	err := Database.DB.Transaction(func(tx *gorm.DB) error {
		// 1. Get the label if exist
		manifestLabel, err := GetManifestLabelByUUID(repositoryId, manifestId, labelUUID)
		if err != nil {
			return err
		}
		if manifestLabel == nil {
			return gorm.ErrRecordNotFound
		}

		// 2. Delete the label association
		err = tx.Delete(manifestLabel).Error
		if err != nil {
			return err
		}

		// 3. Delete the label
		err = tx.Delete(manifestLabel.Label).Error
		if err != nil {
			return err
		}
		return nil
	})
	return err
}
