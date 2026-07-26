package Repositories

import (
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
