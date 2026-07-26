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
